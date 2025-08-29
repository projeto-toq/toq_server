package permissionservice

import (
	"context"
	"database/sql"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// SwitchActiveRole desativa todos os roles do usuário e ativa apenas o especificado
func (p *permissionServiceImpl) SwitchActiveRole(ctx context.Context, userID, newRoleID int64) error {
	if userID <= 0 {
		return utils.ErrBadRequest
	}

	if newRoleID <= 0 {
		return utils.ErrBadRequest
	}

	slog.Debug("Switching active role", "userID", userID, "newRoleID", newRoleID)

	// 1. Desativar todos os roles do usuário
	err := p.DeactivateAllUserRoles(ctx, userID)
	if err != nil {
		return utils.ErrInternalServer
	}

	// 2. Ativar o novo role
	err = p.ActivateUserRole(ctx, userID, newRoleID)
	if err != nil {
		return utils.ErrInternalServer
	}

	slog.Info("Active role switched successfully", "userID", userID, "newRoleID", newRoleID)
	return nil
}

// SwitchActiveRoleWithTx desativa todos os roles do usuário e ativa apenas o especificado (com transação - uso em fluxos)
func (p *permissionServiceImpl) SwitchActiveRoleWithTx(ctx context.Context, tx *sql.Tx, userID, newRoleID int64) error {
	if userID <= 0 {
		return utils.ErrBadRequest
	}

	if newRoleID <= 0 {
		return utils.ErrBadRequest
	}

	slog.Debug("Switching active role with tx", "userID", userID, "newRoleID", newRoleID)

	// 1. Desativar todos os roles do usuário
	err := p.permissionRepository.DeactivateAllUserRoles(ctx, tx, userID)
	if err != nil {
		return utils.ErrInternalServer
	}

	// 2. Ativar o novo role
	err = p.permissionRepository.ActivateUserRole(ctx, tx, userID, newRoleID)
	if err != nil {
		return utils.ErrInternalServer
	}

	slog.Info("Active role switched successfully with tx", "userID", userID, "newRoleID", newRoleID)
	return nil
}

// GetActiveUserRole retorna o role ativo do usuário
func (p *permissionServiceImpl) GetActiveUserRole(ctx context.Context, userID int64) (permissionmodel.UserRoleInterface, error) {
	if userID <= 0 {
		return nil, utils.ErrBadRequest
	}

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		return nil, utils.ErrInternalServer
	}
	defer p.globalService.RollbackTransaction(ctx, tx)

	userRole, err := p.permissionRepository.GetActiveUserRoleByUserID(ctx, tx, userID)
	if err != nil {
		return nil, utils.ErrInternalServer
	}

	// Commit the transaction
	err = p.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		return nil, utils.ErrInternalServer
	}

	return userRole, nil
}

// DeactivateAllUserRoles desativa todos os roles de um usuário
func (p *permissionServiceImpl) DeactivateAllUserRoles(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return utils.ErrBadRequest
	}

	slog.Debug("Deactivating all user roles", "userID", userID)

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		return utils.ErrInternalServer
	}
	defer p.globalService.RollbackTransaction(ctx, tx)

	err = p.permissionRepository.DeactivateAllUserRoles(ctx, tx, userID)
	if err != nil {
		return utils.ErrInternalServer
	}

	// Commit the transaction
	err = p.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		return utils.ErrInternalServer
	}

	slog.Info("All user roles deactivated successfully", "userID", userID)
	return nil
}

// GetRoleBySlug busca um role pelo slug (sem transação - uso direto)
func (p *permissionServiceImpl) GetRoleBySlug(ctx context.Context, slug string) (permissionmodel.RoleInterface, error) {
	if slug == "" {
		return nil, utils.ErrBadRequest
	}

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		return nil, utils.ErrInternalServer
	}
	defer p.globalService.RollbackTransaction(ctx, tx)

	role, err := p.GetRoleBySlugWithTx(ctx, tx, slug)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = p.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		return nil, utils.ErrInternalServer
	}

	return role, nil
}

// GetRoleBySlugWithTx busca um role pelo slug (com transação - uso em fluxos)
func (p *permissionServiceImpl) GetRoleBySlugWithTx(ctx context.Context, tx *sql.Tx, slug string) (permissionmodel.RoleInterface, error) {
	if slug == "" {
		return nil, utils.ErrBadRequest
	}

	role, err := p.permissionRepository.GetRoleBySlug(ctx, tx, slug)
	if err != nil {
		return nil, utils.ErrInternalServer
	}

	if role == nil {
		return nil, utils.ErrNotFound
	}

	return role, nil
}

// ActivateUserRole ativa um role específico do usuário (helper method)
func (p *permissionServiceImpl) ActivateUserRole(ctx context.Context, userID, roleID int64) error {
	if userID <= 0 {
		return utils.ErrBadRequest
	}

	if roleID <= 0 {
		return utils.ErrBadRequest
	}

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		return utils.ErrInternalServer
	}
	defer p.globalService.RollbackTransaction(ctx, tx)

	err = p.permissionRepository.ActivateUserRole(ctx, tx, userID, roleID)
	if err != nil {
		return utils.ErrInternalServer
	}

	// Commit the transaction
	err = p.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		return utils.ErrInternalServer
	}

	return nil
}
