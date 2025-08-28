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

	// TODO: Implementar invalidação de cache quando o cache for implementado
	// Por enquanto, apenas logamos
	slog.Info("User cache invalidated", "userID", userID)
	return nil
}
