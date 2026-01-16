package userservices

import (
	"context"
	"database/sql"
	"strings"
	"time"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

func (us *userService) AddAlternativeRole(ctx context.Context, userID int64, roleSlug permissionmodel.RoleSlug, creciInfo ...string) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("user.add_alternative_role.tracer_error", "err", err)
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		utils.LoggerFromContext(ctx).Error("user.add_alternative_role.tx_start_error", "err", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.LoggerFromContext(ctx).Error("user.add_alternative_role.tx_rollback_error", "err", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	err = us.addAlternativeRole(ctx, tx, userID, roleSlug, creciInfo...)
	if err != nil {
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		utils.LoggerFromContext(ctx).Error("user.add_alternative_role.tx_commit_error", "err", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) addAlternativeRole(ctx context.Context, tx *sql.Tx, userID int64, roleSlug permissionmodel.RoleSlug, creciInfo ...string) (err error) {
	ctx = utils.ContextWithLogger(ctx)

	//verify if the user is on active status
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.add_alternative_role.repo_get_user_by_id_error", "user_id", userID, "err", err)
		return utils.InternalError("Failed to get user")
	}

	// Check if user has active role
	activeRole := user.GetActiveRole()
	if activeRole == nil {
		derr := utils.InternalError("Active role missing")
		utils.LoggerFromContext(ctx).Error("user.active_role.missing", "user_id", userID)
		utils.SetSpanError(ctx, derr)
		return derr
	}

	currentRoleSlug := utils.GetUserRoleSlugFromUserRole(activeRole)
	if currentRoleSlug != permissionmodel.RoleSlugOwner && currentRoleSlug != permissionmodel.RoleSlugRealtor {
		return utils.AuthorizationError("Only owners or realtors can request an alternative role")
	}

	if activeRole.GetStatus() != globalmodel.StatusActive {
		return utils.ConflictError("Active role status must be active")
	}

	expectedAlternative := permissionmodel.RoleSlugOwner
	if currentRoleSlug == permissionmodel.RoleSlugOwner {
		expectedAlternative = permissionmodel.RoleSlugRealtor
	}
	if roleSlug != expectedAlternative {
		return utils.BadRequest("Invalid alternative role for current role")
	}

	var (
		targetStatus  globalmodel.UserRoleStatus
		creciWasSaved bool
	)
	switch roleSlug {
	case permissionmodel.RoleSlugOwner:
		targetStatus = globalmodel.StatusActive
	case permissionmodel.RoleSlugRealtor:
		targetStatus = globalmodel.StatusPendingCreci
	default:
		return utils.AuthorizationError("Unsupported alternative role")
	}

	if roleSlug == permissionmodel.RoleSlugRealtor {
		creciWasSaved, err = us.applyCreciInfoToUser(ctx, tx, user, creciInfo)
		if err != nil {
			return err
		}
	}

	// Get role from permission service
	role, err := us.permissionService.GetRoleBySlugWithTx(ctx, tx, roleSlug)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.add_alternative_role.permission_get_role_error", "user_id", userID, "role", string(roleSlug), "err", err)
		return utils.InternalError("Failed to get role")
	}

	// Create user role using permission service (not active by default)
	isActive := false
	options := &AssignRoleOptions{
		IsActive: &isActive,
		Status:   &targetStatus,
	}

	userRole, err := us.AssignRoleToUserWithTx(ctx, tx, userID, role.GetID(), nil, options)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.add_alternative_role.permission_assign_role_error", "user_id", userID, "role_id", role.GetID(), "err", err)
		return utils.InternalError("Failed to assign role to user")
	}

	// Handle realtor-specific setup
	if roleSlug == permissionmodel.RoleSlugRealtor {
		err = us.CreateUserFolder(ctx, user.GetID())
		if err != nil {
			utils.SetSpanError(ctx, err)
			utils.LoggerFromContext(ctx).Error("user.add_alternative_role.create_user_folder_error", "user_id", user.GetID(), "err", err)
			return utils.InternalError("Failed to create user folder")
		}
	}

	if creciWasSaved {
		auditCreci := auditservice.BuildRecordFromContext(
			ctx,
			userID,
			auditmodel.AuditTarget{Type: auditmodel.TargetUser, ID: userID},
			auditmodel.OperationUpdate,
			map[string]any{
				"action":         "creci_updated_for_alternative_role",
				"role_slug":      roleSlug,
				"creci_number":   user.GetCreciNumber(),
				"creci_state":    user.GetCreciState(),
				"creci_validity": user.GetCreciValidity(),
			},
		)
		if auditErr := us.auditService.RecordChange(ctx, tx, auditCreci); auditErr != nil {
			utils.SetSpanError(ctx, auditErr)
			utils.LoggerFromContext(ctx).Error("user.add_alternative_role.audit_user_creci_error", "table", string(auditmodel.TargetUser), "err", auditErr)
			return utils.InternalError("Failed to create audit record")
		}
	}

	auditRole := auditservice.BuildRecordFromContext(
		ctx,
		userID,
		auditmodel.AuditTarget{Type: auditmodel.TargetUserRole, ID: userRole.GetID()},
		auditmodel.OperationCreate,
		map[string]any{
			"action":        "alternative_role_created",
			"role_slug":     roleSlug,
			"target_status": targetStatus,
			"is_active":     isActive,
		},
	)
	if auditErr := us.auditService.RecordChange(ctx, tx, auditRole); auditErr != nil {
		utils.SetSpanError(ctx, auditErr)
		utils.LoggerFromContext(ctx).Error("user.add_alternative_role.audit_create_error", "table", string(auditmodel.TargetUserRole), "err", auditErr)
		return utils.InternalError("Failed to create audit record")
	}

	return
}

// applyCreciInfoToUser validates and persists CRECI data when assigning realtor role.
func (us *userService) applyCreciInfoToUser(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface, creciInfo []string) (bool, error) {
	logger := utils.LoggerFromContext(ctx)

	if len(creciInfo) == 0 {
		return false, nil
	}

	if len(creciInfo) != 3 {
		derr := utils.ValidationError("creciInfo", "Realtor role requires full CRECI info (number, state, validity)")
		utils.SetSpanError(ctx, derr)
		return false, derr
	}

	normalizedNumber, err := validators.ValidateCreciNumber("creciNumber", creciInfo[0], true)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return false, err
	}

	normalizedState, err := validators.ValidateCreciState("creciState", creciInfo[1], true)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return false, err
	}

	creciValidityRaw := strings.TrimSpace(creciInfo[2])
	if creciValidityRaw == "" {
		derr := utils.ValidationError("creciValidity", "Creci validity date is required")
		utils.SetSpanError(ctx, derr)
		return false, derr
	}

	creciValidity, perr := time.Parse("2006-01-02", creciValidityRaw)
	if perr != nil {
		derr := utils.ValidationError("creciValidity", "Invalid date format, expected YYYY-MM-DD")
		utils.SetSpanError(ctx, derr)
		return false, derr
	}

	currentNumber := user.GetCreciNumber()
	currentState := user.GetCreciState()
	currentValidity := user.GetCreciValidity()

	// Log when overwriting existing data with different values.
	if currentNumber != "" && currentNumber != normalizedNumber {
		logger.Warn("user.add_alternative_role.creci_number_overwrite", "user_id", user.GetID(), "previous", currentNumber, "new", normalizedNumber)
	}
	if currentState != "" && currentState != normalizedState {
		logger.Warn("user.add_alternative_role.creci_state_overwrite", "user_id", user.GetID(), "previous", currentState, "new", normalizedState)
	}
	if !currentValidity.IsZero() && !currentValidity.Equal(creciValidity) {
		logger.Warn("user.add_alternative_role.creci_validity_overwrite", "user_id", user.GetID(), "previous", currentValidity.String(), "new", creciValidity.String())
	}

	hasChanges := currentNumber != normalizedNumber || currentState != normalizedState || !currentValidity.Equal(creciValidity)
	if !hasChanges {
		return false, nil
	}

	user.SetCreciNumber(normalizedNumber)
	user.SetCreciState(normalizedState)
	user.SetCreciValidity(creciValidity)

	if err := us.repo.UpdateUserByID(ctx, tx, user); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.add_alternative_role.update_creci_error", "user_id", user.GetID(), "err", err)
		return false, utils.InternalError("Failed to update user with CRECI info")
	}

	logger.Info("user.add_alternative_role.creci_updated", "user_id", user.GetID())

	return true, nil
}
