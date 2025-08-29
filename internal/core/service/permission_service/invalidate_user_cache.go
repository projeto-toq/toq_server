package permissionservice

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// InvalidateUserCache invalida o cache de permissões de um usuário
func (p *permissionServiceImpl) InvalidateUserCache(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return utils.ErrBadRequest
	}

	slog.Debug("Invalidating user cache", "userID", userID)

	if p.cache != nil {
		err := p.cache.DeleteUserPermissions(ctx, userID)
		if err != nil {
			slog.Warn("Failed to invalidate user cache in Redis", "userID", userID, "error", err)
			return utils.ErrInternalServer
		}
		slog.Info("User cache invalidated successfully", "userID", userID)
	} else {
		slog.Debug("Cache not available, skipping invalidation", "userID", userID)
	}

	return nil
}
