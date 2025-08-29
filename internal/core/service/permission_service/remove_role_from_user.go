package permissionservice

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// RemoveRoleFromUser remove um role de um usuário
func (p *permissionServiceImpl) RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error {
	if userID <= 0 {
		return utils.ErrBadRequest
	}

	if roleID <= 0 {
		return utils.ErrBadRequest
	}

	slog.Debug("Removing role from user", "userID", userID, "roleID", roleID)

	// Buscar o UserRole existente
	userRole, err := p.permissionRepository.GetUserRoleByUserIDAndRoleID(ctx, nil, userID, roleID)
	if err != nil {
		return utils.ErrInternalServer
	}
	if userRole == nil {
		return utils.ErrNotFound
	}

	// Remover o UserRole
	err = p.permissionRepository.DeleteUserRole(ctx, nil, userRole.GetID())
	if err != nil {
		return utils.ErrInternalServer
	}

	slog.Info("Role removed from user successfully", "userID", userID, "roleID", roleID)
	return nil
}
