package permissionservice

import (
	"context"
	"database/sql"
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// AssignRoleToUser atribui um role a um usuário (sem transação - uso direto)
func (p *permissionServiceImpl) AssignRoleToUser(ctx context.Context, userID, roleID int64, expiresAt *time.Time) (permissionmodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.role.assign.tx_start_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rollbackErr := p.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
				logger.Error("permission.role.assign.tx_rollback_failed", "user_id", userID, "role_id", roleID, "error", rollbackErr)
				utils.SetSpanError(ctx, rollbackErr)
			}
		}
	}()

	userRole, err := p.AssignRoleToUserWithTx(ctx, tx, userID, roleID, expiresAt)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = p.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		logger.Error("permission.role.assign.tx_commit_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return userRole, nil
}

// AssignRoleToUserWithTx atribui um role a um usuário (com transação - uso em fluxos)
func (p *permissionServiceImpl) AssignRoleToUserWithTx(ctx context.Context, tx *sql.Tx, userID, roleID int64, expiresAt *time.Time) (permissionmodel.UserRoleInterface, error) {
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
	role, err := p.permissionRepository.GetRoleByID(ctx, tx, roleID)
	if err != nil {
		logger.Error("permission.role.assign.db_failed", "stage", "get_role", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	if role == nil {
		return nil, utils.NotFoundError("role")
	}

	// Verificar se o usuário já tem este role
	existingUserRole, err := p.permissionRepository.GetUserRoleByUserIDAndRoleID(ctx, tx, userID, roleID)
	if err != nil {
		logger.Error("permission.role.assign.db_failed", "stage", "get_user_role", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	if existingUserRole != nil {
		return nil, utils.ConflictError("role already assigned to user")
	}

	// Criar o novo UserRole
	userRole := permissionmodel.NewUserRole()
	userRole.SetUserID(userID)
	userRole.SetRoleID(roleID)
	userRole.SetIsActive(true)
	userRole.SetStatus(permissionmodel.StatusPendingBoth)

	if expiresAt != nil {
		userRole.SetExpiresAt(expiresAt)
	}

	// Salvar no banco
	userRole, err = p.permissionRepository.CreateUserRole(ctx, tx, userRole)
	if err != nil {
		logger.Error("permission.role.assign.db_failed", "stage", "create_user_role", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	logger.Info("permission.role.assigned", "user_id", userID, "role_id", roleID, "role_name", role.GetName())
	p.invalidateUserCacheSafe(ctx, userID, "assign_role_to_user")
	return userRole, nil
}
