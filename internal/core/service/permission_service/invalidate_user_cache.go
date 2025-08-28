package permissionservice

import (
	"context"
	"fmt"
	"log/slog"
)

// InvalidateUserCache invalida o cache de permissões de um usuário
func (p *permissionServiceImpl) InvalidateUserCache(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	slog.Debug("Invalidating user cache", "userID", userID)

	if p.cache != nil {
		err := p.cache.DeleteUserPermissions(ctx, userID)
		if err != nil {
			slog.Warn("Failed to invalidate user cache in Redis", "userID", userID, "error", err)
			return fmt.Errorf("failed to invalidate user cache: %w", err)
		}
		slog.Info("User cache invalidated successfully", "userID", userID)
	} else {
		slog.Debug("Cache not available, skipping invalidation", "userID", userID)
	}

	return nil
}
