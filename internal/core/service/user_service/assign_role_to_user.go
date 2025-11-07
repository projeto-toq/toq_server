package userservices

import (
	"context"
	"database/sql"
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// AssignRoleOptions permite personalizar campos do user_role criado pelo serviço.
type AssignRoleOptions struct {
	IsActive *bool
	Status   *globalmodel.UserRoleStatus
}

// AssignRoleToUser atribui um role a um usuário (sem transação - uso direto)
func (us *userService) AssignRoleToUser(ctx context.Context, userID, roleID int64, expiresAt *time.Time, opts *AssignRoleOptions) (usermodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.role.assign.tx_start_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rollbackErr := us.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
				logger.Error("permission.role.assign.tx_rollback_failed", "user_id", userID, "role_id", roleID, "error", rollbackErr)
				utils.SetSpanError(ctx, rollbackErr)
			}
		}
	}()

	userRole, err := us.AssignRoleToUserWithTx(ctx, tx, userID, roleID, expiresAt, opts)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		logger.Error("permission.role.assign.tx_commit_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return userRole, nil
}

// AssignRoleToUserWithTx atribui um role a um usuário (com transação - uso em fluxos)
func (us *userService) AssignRoleToUserWithTx(ctx context.Context, tx *sql.Tx, userID, roleID int64, expiresAt *time.Time, opts *AssignRoleOptions) (usermodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	if roleID <= 0 {
		return nil, utils.BadRequest("invalid role id")
	}

	logger.Debug("permission.role.assign.request", "user_id", userID, "role_id", roleID, "expires_at", expiresAt)

	// Verificar se o role existe
	role, err := us.permissionService.GetRoleByIDWithTx(ctx, tx, roleID)
	if err != nil {
		logger.Error("permission.role.assign.db_failed", "stage", "get_role", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	if role == nil {
		return nil, utils.NotFoundError("role")
	}

	// Verificar se o usuário já tem este role
	existingUserRole, err := us.repo.GetUserRoleByUserIDAndRoleID(ctx, tx, userID, roleID)
	if err != nil {
		logger.Error("permission.role.assign.db_failed", "stage", "get_user_role", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	if existingUserRole != nil {
		return nil, utils.ConflictError("role already assigned to user")
	}

	// Criar o novo UserRole
	userRole := usermodel.NewUserRole()
	userRole.SetUserID(userID)
	userRole.SetRoleID(roleID)

	isActive := true
	if opts != nil && opts.IsActive != nil {
		isActive = *opts.IsActive
	}
	userRole.SetIsActive(isActive)

	status := globalmodel.StatusPendingBoth
	if opts != nil && opts.Status != nil {
		status = *opts.Status
	}
	userRole.SetStatus(status)

	if expiresAt != nil {
		userRole.SetExpiresAt(expiresAt)
	}

	// Salvar no banco
	userRole, err = us.repo.CreateUserRole(ctx, tx, userRole)
	if err != nil {
		logger.Error("permission.role.assign.db_failed", "stage", "create_user_role", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	logger.Info("permission.role.assigned", "user_id", userID, "role_id", roleID, "role_name", role.GetName(), "is_active", isActive, "status", status.String())
	us.permissionService.InvalidateUserCache(ctx, userID)
	return userRole, nil
}
