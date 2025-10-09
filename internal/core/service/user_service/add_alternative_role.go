package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
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

	// Validate creci info for realtor role
	if roleSlug == permissionmodel.RoleSlugRealtor && len(creciInfo) != 3 {
		return utils.ValidationError("creciInfo", "Realtor role requires CRECI info")
	}

	// Get role from permission service
	role, err := us.permissionService.GetRoleBySlugWithTx(ctx, tx, roleSlug)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.add_alternative_role.permission_get_role_error", "user_id", userID, "role", string(roleSlug), "err", err)
		return utils.InternalError("Failed to get role")
	}

	// Create user role using permission service (not active by default)
	_, err = us.permissionService.AssignRoleToUserWithTx(ctx, tx, userID, role.GetID(), nil)
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

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUserRoles, "Criado papel alternativo")
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.add_alternative_role.audit_create_error", "table", string(globalmodel.TableUserRoles), "err", err)
		return utils.InternalError("Failed to create audit record")
	}

	return
}
