package permissionservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// InvalidateUserCache invalida o cache de permissões de um usuário
//
// Esta função remove todas as chaves de cache relacionadas ao usuário:
//  1. Cache agregado de permissões (DeleteUserPermissions)
//  2. Chaves granulares por action (CleanByUser - best-effort)
//
// Registra métricas de cache com label 'source' para rastreabilidade.
//
// Parameters:
//   - ctx: Context for tracing and logging (must contain request metadata)
//   - userID: ID do usuário cujo cache será invalidado (must be > 0)
//   - source: Identificador da operação que causou a invalidação
//     Exemplos: "assign_role", "remove_role", "switch_active_role", "grant_permission_to_role"
//     Usado em logs e métricas para rastreabilidade
//
// Returns:
//   - error: Infrastructure error se falhar ao invalidar o cache agregado
//
// Business Rules:
//   - Se cache não estiver configurado (nil), retorna sucesso sem operação
//   - Falha ao limpar chaves granulares NÃO retorna erro (best-effort)
//   - Registra métricas Prometheus com label 'operation' = 'user_permissions_invalidate'
//
// Example:
//
//	// Após atribuir role ao usuário
//	if err := ps.InvalidateUserCache(ctx, userID, "assign_role"); err != nil {
//	    logger.Error("cache_invalidation_failed", "error", err)
//	    return err
//	}
func (p *permissionServiceImpl) InvalidateUserCache(ctx context.Context, userID int64, source string) error {
	// Initialize tracing for distributed observability
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	// Ensure logger propagation with request_id and trace_id
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Validate user ID (business rule)
	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	logger.Debug("permission.cache.invalidate.request", "user_id", userID, "source", source)

	// Check if cache is enabled (infrastructure availability)
	if p.cache == nil {
		logger.Debug("permission.cache.not_available", "user_id", userID, "source", source)
		p.observeCacheOperation("user_permissions_invalidate", "disabled")
		return nil
	}

	// Clear aggregated user permissions cache (critical operation)
	if err := p.ClearUserPermissionsCache(ctx, userID); err != nil {
		// ClearUserPermissionsCache already logged the error; add context
		logger.Warn("permission.cache.clear_failed", "user_id", userID, "source", source, "error", err)
		p.observeCacheOperation("user_permissions_invalidate", "error")
		return err
	}

	// Clear granular permission keys (best-effort, non-critical)
	// Does not return error to avoid blocking the operation
	p.cache.CleanByUser(ctx, userID)

	logger.Info("permission.cache.invalidated", "user_id", userID, "source", source)
	p.observeCacheOperation("user_permissions_invalidate", "success")
	return nil
}

// InvalidateUserCacheSafe tenta invalidar o cache e apenas registra aviso em caso de falha
//
// Este método é ideal para operações best-effort, onde a falha na invalidação
// não deve bloquear o fluxo principal (ex: após commit de transação).
//
// NÃO retorna erro. Em caso de falha:
//   - Marca span como erro para tracing
//   - Registra WARN em logs com contexto completo
//   - Não interrompe a operação do caller
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - userID: ID do usuário cujo cache será invalidado (must be > 0)
//   - source: Identificador da operação que causou a invalidação
//
// Usage:
//
//	// Após commit bem-sucedido, quando falha não deve reverter a operação
//	ps.InvalidateUserCacheSafe(ctx, userID, "assign_role")
func (p *permissionServiceImpl) InvalidateUserCacheSafe(ctx context.Context, userID int64, source string) {
	// Ensure logger propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Skip if user ID is invalid (silent fail for best-effort)
	if userID <= 0 {
		return
	}

	// Attempt invalidation; log warning on failure but don't propagate error
	if err := p.InvalidateUserCache(ctx, userID, source); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Warn("permission.cache.invalidate_safe_failed", "user_id", userID, "source", source, "error", err)
	}
}
