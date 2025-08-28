package permissionservice

import (
	"context"
	"fmt"
	"log/slog"
)

// RemoveRoleFromUser remove um role de um usuário
func (p *permissionServiceImpl) RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	if roleID <= 0 {
		return fmt.Errorf("invalid role ID: %d", roleID)
	}

	slog.Debug("Removing role from user", "userID", userID, "roleID", roleID)

	// Buscar o UserRole existente
	userRole, err := p.permissionRepository.GetUserRoleByUserIDAndRoleID(ctx, nil, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to find user role: %w", err)
	}
	if userRole == nil {
		return fmt.Errorf("user %d does not have role %d", userID, roleID)
	}

	// Remover o UserRole
	err = p.permissionRepository.DeleteUserRole(ctx, nil, userRole.GetID())
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	slog.Info("Role removed from user successfully", "userID", userID, "roleID", roleID)
	return nil
}
