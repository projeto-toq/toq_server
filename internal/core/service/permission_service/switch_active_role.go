package permissionservice

import (
	"context"
	"database/sql"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// SwitchActiveRole desativa todos os roles do usuário e ativa apenas o especificado
func (p *permissionServiceImpl) SwitchActiveRole(ctx context.Context, userID, newRoleID int64) error {
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
	if err := p.DeactivateAllUserRoles(ctx, userID); err != nil {
		logger.Error("permission.role.switch.deactivate_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return err
	}

	// 2. Ativar o novo role
	if err := p.ActivateUserRole(ctx, userID, newRoleID); err != nil {
		logger.Error("permission.role.switch.activate_failed", "user_id", userID, "new_role_id", newRoleID, "error", err)
		utils.SetSpanError(ctx, err)
		return err
	}

	logger.Info("permission.role.switched", "user_id", userID, "new_role_id", newRoleID)
	return nil
}

// SwitchActiveRoleWithTx desativa todos os roles do usuário e ativa apenas o especificado (com transação - uso em fluxos)
func (p *permissionServiceImpl) SwitchActiveRoleWithTx(ctx context.Context, tx *sql.Tx, userID, newRoleID int64) error {
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
	if err := p.permissionRepository.DeactivateAllUserRoles(ctx, tx, userID); err != nil {
		logger.Error("permission.role.switch.tx.deactivate_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// 2. Ativar o novo role
	if err := p.permissionRepository.ActivateUserRole(ctx, tx, userID, newRoleID); err != nil {
		logger.Error("permission.role.switch.tx.activate_failed", "user_id", userID, "new_role_id", newRoleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	logger.Info("permission.role.switched.tx", "user_id", userID, "new_role_id", newRoleID)
	p.invalidateUserCacheSafe(ctx, userID, "switch_active_role_with_tx")
	return nil
}

// GetActiveUserRole methods moved to get_active_user_role.go

// DeactivateAllUserRoles desativa todos os roles de um usuário
func (p *permissionServiceImpl) DeactivateAllUserRoles(ctx context.Context, userID int64) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	logger.Debug("permission.user_roles.deactivate.request", "user_id", userID)

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user_roles.deactivate.tx_start_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user_roles.deactivate.tx_rollback_failed", "user_id", userID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	if err = p.permissionRepository.DeactivateAllUserRoles(ctx, tx, userID); err != nil {
		logger.Error("permission.user_roles.deactivate.db_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Commit the transaction
	if err = p.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.user_roles.deactivate.tx_commit_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	logger.Info("permission.user_roles.deactivated", "user_id", userID)
	p.invalidateUserCacheSafe(ctx, userID, "deactivate_all_user_roles")
	return nil
}

// GetRoleBySlug busca um role pelo slug (sem transação - uso direto)
func (p *permissionServiceImpl) GetRoleBySlug(ctx context.Context, slug permissionmodel.RoleSlug) (permissionmodel.RoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if slug == "" {
		return nil, utils.BadRequest("invalid role slug")
	}

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.role.get_by_slug.tx_start_failed", "slug", slug, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.role.get_by_slug.tx_rollback_failed", "slug", slug, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	role, derr := p.GetRoleBySlugWithTx(ctx, tx, slug)
	if derr != nil {
		return nil, derr
	}

	// Commit the transaction
	if err = p.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.role.get_by_slug.tx_commit_failed", "slug", slug, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return role, nil
}

// GetRoleBySlugWithTx busca um role pelo slug (com transação - uso em fluxos)
func (p *permissionServiceImpl) GetRoleBySlugWithTx(ctx context.Context, tx *sql.Tx, slug permissionmodel.RoleSlug) (permissionmodel.RoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	role, err := p.permissionRepository.GetRoleBySlug(ctx, tx, slug.String())
	if err != nil {
		logger.Error("permission.role.get_by_slug.db_failed", "slug", slug, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return role, nil
}

// ActivateUserRole ativa um role específico do usuário (helper method)
func (p *permissionServiceImpl) ActivateUserRole(ctx context.Context, userID, roleID int64) error {
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
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user_role.activate.tx_start_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user_role.activate.tx_rollback_failed", "user_id", userID, "role_id", roleID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	if err = p.permissionRepository.ActivateUserRole(ctx, tx, userID, roleID); err != nil {
		logger.Error("permission.user_role.activate.db_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Commit the transaction
	if err = p.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.user_role.activate.tx_commit_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	p.invalidateUserCacheSafe(ctx, userID, "activate_user_role")
	return nil
}
