package userservices

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// RemoveRoleFromUser remove um role de um usuário
func (us *userService) RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error {
	// Tracing da operação
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	if roleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	logger.Debug("permission.role.remove.start", "user_id", userID, "role_id", roleID)

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.role.remove.tx_start_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.role.remove.tx_rollback_failed", "user_id", userID, "role_id", roleID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	if err = us.RemoveRoleFromUserWithTx(ctx, tx, userID, roleID); err != nil {
		return err
	}

	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.role.remove.tx_commit_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	return nil
}

// RemoveRoleFromUserWithTx remove um role de um usuário usando uma transação existente
func (us *userService) RemoveRoleFromUserWithTx(ctx context.Context, tx *sql.Tx, userID, roleID int64) error {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	if roleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	logger.Debug("permission.role.remove.start.tx", "user_id", userID, "role_id", roleID)

	// Buscar o UserRole existente
	userRole, err := us.repo.GetUserRoleByUserIDAndRoleID(ctx, tx, userID, roleID)
	if err != nil {
		logger.Error("permission.role.get_user_role_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	if userRole == nil {
		return utils.NotFoundError("UserRole")
	}

	// Remover o UserRole
	err = us.repo.DeleteUserRole(ctx, tx, userRole.GetID())
	if err != nil {
		logger.Error("permission.role.delete_user_role_failed", "user_id", userID, "role_id", roleID, "user_role_id", userRole.GetID(), "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	logger.Info("permission.role.removed", "user_id", userID, "role_id", roleID)
	us.permissionService.InvalidateUserCache(ctx, userID)
	return nil
}
