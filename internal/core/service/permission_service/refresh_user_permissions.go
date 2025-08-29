package permissionservice

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// RefreshUserPermissions atualiza o cache de permissões de um usuário
func (p *permissionServiceImpl) RefreshUserPermissions(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return utils.ErrBadRequest
	}

	slog.Debug("Refreshing user permissions", "userID", userID)

	// Invalidar cache atual
	if err := p.InvalidateUserCache(ctx, userID); err != nil {
		slog.Warn("Failed to invalidate cache before refresh", "userID", userID, "error", err)
	}

	// Buscar permissões atuais do banco e recriar cache
	_, err := p.getUserPermissionsWithCache(ctx, userID)
	if err != nil {
		return utils.ErrInternalServer
	}

	slog.Info("User permissions refreshed", "userID", userID)
	return nil
}
