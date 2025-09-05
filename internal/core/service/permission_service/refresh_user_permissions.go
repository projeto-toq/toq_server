package permissionservice

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// RefreshUserPermissions atualiza o cache de permissões de um usuário
func (p *permissionServiceImpl) RefreshUserPermissions(ctx context.Context, userID int64) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	slog.Debug("permission.permissions.refresh.request", "user_id", userID)

	// Invalidar cache atual
	if err := p.InvalidateUserCache(ctx, userID); err != nil {
		slog.Warn("permission.permissions.refresh.invalidate_failed", "user_id", userID, "error", err)
	}

	// Buscar permissões atuais do banco e recriar cache
	_, err := p.getUserPermissionsWithCache(ctx, userID)
	if err != nil {
		slog.Error("permission.permissions.refresh.db_or_cache_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	slog.Info("permission.permissions.refreshed", "user_id", userID)
	return nil
}
