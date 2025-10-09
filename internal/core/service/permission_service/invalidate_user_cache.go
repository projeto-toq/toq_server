package permissionservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// InvalidateUserCache invalida o cache de permissões de um usuário
func (p *permissionServiceImpl) InvalidateUserCache(ctx context.Context, userID int64) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	logger.Debug("permission.cache.invalidate.request", "user_id", userID)

	// Se o cache não estiver configurado, apenas registrar e sair
	if p.cache == nil {
		logger.Debug("permission.cache.not_available", "user_id", userID)
		p.observeCacheOperation("user_permissions_invalidate", "disabled")
		return nil
	}

	// Primeiro, remover o cache de permissões agregadas do usuário
	if err := p.ClearUserPermissionsCache(ctx, userID); err != nil {
		// ClearUserPermissionsCache já padroniza o erro; apenas logamos contexto
		logger.Warn("permission.cache.clear_failed", "user_id", userID, "error", err)
		p.observeCacheOperation("user_permissions_invalidate", "error")
		return err
	}

	// Em seguida, limpar chaves granulares de permissões (resource/action)
	// Best-effort granular clean (sem retorno de erro)
	p.cache.CleanByUser(ctx, userID)

	logger.Info("permission.cache.invalidated", "user_id", userID)
	p.observeCacheOperation("user_permissions_invalidate", "success")
	return nil
}

// invalidateUserCacheSafe tenta invalidar o cache e apenas registra aviso em caso de falha
func (p *permissionServiceImpl) invalidateUserCacheSafe(ctx context.Context, userID int64, source string) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	if userID <= 0 {
		return
	}

	if err := p.InvalidateUserCache(ctx, userID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Warn("permission.cache.invalidate_safe_failed", "user_id", userID, "source", source, "error", err)
	}
}
