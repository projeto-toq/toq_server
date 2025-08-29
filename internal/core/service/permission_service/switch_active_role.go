package permissionservice

import (
	"context"
	"fmt"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// SwitchActiveRole desativa todos os roles do usuário e ativa apenas o especificado
func (p *permissionServiceImpl) SwitchActiveRole(ctx context.Context, userID, newRoleID int64) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	if newRoleID <= 0 {
		return fmt.Errorf("invalid role ID: %d", newRoleID)
	}

	slog.Debug("Switching active role", "userID", userID, "newRoleID", newRoleID)

	// 1. Desativar todos os roles do usuário
	err := p.DeactivateAllUserRoles(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate all user roles: %w", err)
	}

	// 2. Ativar o novo role
	err = p.ActivateUserRole(ctx, userID, newRoleID)
	if err != nil {
		return fmt.Errorf("failed to activate new role: %w", err)
	}

	slog.Info("Active role switched successfully", "userID", userID, "newRoleID", newRoleID)
	return nil
}

// GetActiveUserRole retorna o role ativo do usuário
func (p *permissionServiceImpl) GetActiveUserRole(ctx context.Context, userID int64) (permissionmodel.UserRoleInterface, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	userRole, err := p.permissionRepository.GetActiveUserRoleByUserID(ctx, nil, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active user role: %w", err)
	}

	return userRole, nil
}

// DeactivateAllUserRoles desativa todos os roles de um usuário
func (p *permissionServiceImpl) DeactivateAllUserRoles(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	slog.Debug("Deactivating all user roles", "userID", userID)

	err := p.permissionRepository.DeactivateAllUserRoles(ctx, nil, userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate all user roles: %w", err)
	}

	slog.Info("All user roles deactivated successfully", "userID", userID)
	return nil
}

// GetRoleBySlug busca um role pelo slug
func (p *permissionServiceImpl) GetRoleBySlug(ctx context.Context, slug string) (permissionmodel.RoleInterface, error) {
	if slug == "" {
		return nil, fmt.Errorf("slug cannot be empty")
	}

	role, err := p.permissionRepository.GetRoleBySlug(ctx, nil, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by slug: %w", err)
	}

	if role == nil {
		return nil, fmt.Errorf("role with slug '%s' not found", slug)
	}

	return role, nil
}

// ActivateUserRole ativa um role específico do usuário (helper method)
func (p *permissionServiceImpl) ActivateUserRole(ctx context.Context, userID, roleID int64) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	if roleID <= 0 {
		return fmt.Errorf("invalid role ID: %d", roleID)
	}

	err := p.permissionRepository.ActivateUserRole(ctx, nil, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to activate user role: %w", err)
	}

	return nil
}
