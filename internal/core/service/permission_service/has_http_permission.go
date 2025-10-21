package permissionservice

import (
	"context"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// HasHTTPPermission verifica se o usuário tem permissão para um endpoint HTTP específico
func (p *permissionServiceImpl) HasHTTPPermission(ctx context.Context, userID int64, method, path string) (bool, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return false, utils.BadRequest("invalid user id")
	}

	if method == "" || path == "" {
		return false, utils.BadRequest("invalid http method or path")
	}

	action := fmt.Sprintf("%s:%s", method, path)
	logger.Debug("permission.http.check.start", "user_id", userID, "action", action)

	allowed, err := p.hasPermissionByAction(ctx, userID, action)
	if err != nil {
		logger.Error("permission.http.check.error", "user_id", userID, "action", action, "error", err)
		return false, utils.InternalError("")
	}

	if !allowed {
		logger.Warn("permission.http.check.denied", "user_id", userID, "action", action)
	}

	return allowed, nil
}
