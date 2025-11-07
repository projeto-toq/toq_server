package userservices

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// SwitchActiveRole desativa todos os roles do usuário e ativa apenas o especificado
func (us *userService) SwitchActiveRole(ctx context.Context, userID, newRoleID int64) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	if newRoleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	logger.Debug("permission.role.switch.request", "user_id", userID, "new_role_id", newRoleID)

	// 1. Desativar todos os roles do usuário
	if err := us.DeactivateAllUserRoles(ctx, userID); err != nil {
		logger.Error("permission.role.switch.deactivate_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return err
	}

	// 2. Ativar o novo role
	if err := us.ActivateUserRole(ctx, userID, newRoleID); err != nil {
		logger.Error("permission.role.switch.activate_failed", "user_id", userID, "new_role_id", newRoleID, "error", err)
		utils.SetSpanError(ctx, err)
		return err
	}

	logger.Info("permission.role.switched", "user_id", userID, "new_role_id", newRoleID)
	return nil
}

// SwitchActiveRoleWithTx desativa todos os roles do usuário e ativa apenas o especificado (com transação - uso em fluxos)
func (us *userService) SwitchActiveRoleWithTx(ctx context.Context, tx *sql.Tx, userID, newRoleID int64) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	if newRoleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	logger.Debug("permission.role.switch.tx.request", "user_id", userID, "new_role_id", newRoleID)

	// 1. Desativar todos os roles do usuário
	if err := us.repo.DeactivateAllUserRoles(ctx, tx, userID); err != nil {
		logger.Error("permission.role.switch.tx.deactivate_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// 2. Ativar o novo role
	if err := us.repo.ActivateUserRole(ctx, tx, userID, newRoleID); err != nil {
		logger.Error("permission.role.switch.tx.activate_failed", "user_id", userID, "new_role_id", newRoleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	logger.Info("permission.role.switched.tx", "user_id", userID, "new_role_id", newRoleID)
	us.permissionService.InvalidateUserCache(ctx, userID) // TODO incluir mensagem, "switch_active_role_with_tx")
	return nil
}

// GetActiveUserRole methods moved to get_active_user_role.go

// DeactivateAllUserRoles desativa todos os roles de um usuário
func (us *userService) DeactivateAllUserRoles(ctx context.Context, userID int64) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	logger.Debug("permission.user_roles.deactivate.request", "user_id", userID)

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user_roles.deactivate.tx_start_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user_roles.deactivate.tx_rollback_failed", "user_id", userID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	if err = us.repo.DeactivateAllUserRoles(ctx, tx, userID); err != nil {
		logger.Error("permission.user_roles.deactivate.db_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Commit the transaction
	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.user_roles.deactivate.tx_commit_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	logger.Info("permission.user_roles.deactivated", "user_id", userID)
	us.permissionService.InvalidateUserCache(ctx, userID) //TODO incluir mensatgem, "deactivate_all_user_roles")
	return nil
}

// ActivateUserRole ativa um role específico do usuário (helper method)
func (us *userService) ActivateUserRole(ctx context.Context, userID, roleID int64) error {
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

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user_role.activate.tx_start_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user_role.activate.tx_rollback_failed", "user_id", userID, "role_id", roleID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	if err = us.repo.ActivateUserRole(ctx, tx, userID, roleID); err != nil {
		logger.Error("permission.user_role.activate.db_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Commit the transaction
	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.user_role.activate.tx_commit_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	us.permissionService.InvalidateUserCache(ctx, userID) //TODO necessário ajustar para colocar msg, "activate_user_role")
	return nil
}
