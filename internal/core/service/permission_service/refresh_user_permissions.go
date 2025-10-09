package permissionservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// RefreshUserPermissions atualiza o cache de permissões de um usuário
func (p *permissionServiceImpl) RefreshUserPermissions(ctx context.Context, userID int64) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	logger.Debug("permission.permissions.refresh.request", "user_id", userID)

	// Invalidar cache atual
	if err := p.InvalidateUserCache(ctx, userID); err != nil {
		logger.Warn("permission.permissions.refresh.invalidate_failed", "user_id", userID, "error", err)
	}

	// Buscar permissões atuais do banco e recriar cache
	_, _, err := p.getUserPermissionsWithCache(ctx, userID)
	if err != nil {
		logger.Error("permission.permissions.refresh.db_or_cache_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		p.observeCacheOperation("user_permissions_refresh", "error")
		return utils.InternalError("")
	}

	logger.Info("permission.permissions.refreshed", "user_id", userID)
	p.observeCacheOperation("user_permissions_refresh", "success")
	return nil
}
