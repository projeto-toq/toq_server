package middlewares

import (
	"github.com/gin-gonic/gin"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	permissionservice "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PermissionMiddleware verifica permissões específicas usando o sistema de permissões avançado
func PermissionMiddleware(permissionService permissionservice.PermissionServiceInterface) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		method := c.Request.Method
		path := c.Request.URL.Path

		// Skip para endpoints públicos
		if !isPermissionCheckRequired(method, path) {
			c.Next()
			return
		}

		ctx := c.Request.Context()
		logger := coreutils.LoggerFromContext(ctx)

		// Obter informações do usuário do contexto (definido pelo auth middleware)
		userInfoInterface, exists := c.Get("userInfo")
		if !exists {
			httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("User info not found"))
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("permission", "missing_user_info")
			}
			logger.Warn("permission.middleware.user_info_missing", "path", path, "method", method)
			c.Abort()
			return
		}

		userInfo, ok := userInfoInterface.(usermodel.UserInfos)
		if !ok {
			httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("Invalid user info format"))
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("permission", "invalid_user_info")
			}
			logger.Warn("permission.middleware.user_info_invalid", "path", path, "method", method)
			c.Abort()
			return
		}

		userID := userInfo.ID

		// Usar o novo sistema de permissões HTTP
		hasPermission, err := permissionService.HasHTTPPermission(ctx, userID, method, path)
		if err != nil {
			logger.Error("permission.middleware.check_failed", "user_id", userID, "method", method, "path", path, "err", err)
			httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Permission check failed"))
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("permission", "check_failed")
			}
			c.Abort()
			return
		}

		if !hasPermission {
			logger.Warn("permission.middleware.denied", "user_id", userID, "method", method, "path", path)
			httperrors.SendHTTPErrorObj(c, coreutils.AuthorizationError("Insufficient permissions"))
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("permission", "forbidden")
			}
			c.Abort()
			return
		}

		logger.Debug("permission.middleware.granted", "user_id", userID, "method", method, "path", path)
		c.Next()
	})
}

// isPermissionCheckRequired verifica se o endpoint precisa de verificação de permissão
func isPermissionCheckRequired(method, path string) bool {
	// Use a single source of truth for public endpoints
	// Method param kept for possible future expansion (e.g., method-specific rules)
	_ = method
	return !coreutils.IsPublicEndpoint(path)
}
