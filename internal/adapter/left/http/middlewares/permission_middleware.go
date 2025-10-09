package middlewares

import (
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	permissionservice "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// PermissionMiddleware verifica permissões específicas usando o sistema de permissões avançado
func PermissionMiddleware(permissionService permissionservice.PermissionServiceInterface) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Skip para endpoints públicos
		if !isPermissionCheckRequired(c.Request.Method, c.Request.URL.Path) {
			c.Next()
			return
		}

		// Obter informações do usuário do contexto (definido pelo auth middleware)
		userInfoInterface, exists := c.Get("userInfo")
		if !exists {
			httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("User info not found"))
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("permission", "missing_user_info")
			}
			c.Abort()
			return
		}

		userInfo, ok := userInfoInterface.(usermodel.UserInfos)
		if !ok {
			httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("Invalid user info format"))
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("permission", "invalid_user_info")
			}
			c.Abort()
			return
		}

		userID := userInfo.ID
		method := c.Request.Method
		path := c.Request.URL.Path

		// Usar o novo sistema de permissões HTTP
		hasPermission, err := permissionService.HasHTTPPermission(c.Request.Context(), userID, method, path)
		if err != nil {
			slog.Error("Error checking permission", "userID", userID, "method", method, "path", path, "error", err)
			httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Permission check failed"))
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("permission", "check_failed")
			}
			c.Abort()
			return
		}

		if !hasPermission {
			slog.Warn("Permission denied", "userID", userID, "method", method, "path", path)
			httperrors.SendHTTPErrorObj(c, coreutils.AuthorizationError("Insufficient permissions"))
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("permission", "forbidden")
			}
			c.Abort()
			return
		}

		slog.Debug("Permission granted", "userID", userID, "method", method, "path", path)
		c.Next()
	})
}

// RequirePermission middleware para verificar permissões específicas por resource/action
func RequirePermission(permissionService permissionservice.PermissionServiceInterface, resource, action string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Obter informações do usuário
		userInfoInterface, exists := c.Get("userInfo")
		if !exists {
			httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("User info not found"))
			c.Abort()
			return
		}

		userInfo, ok := userInfoInterface.(usermodel.UserInfos)
		if !ok {
			httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("Invalid user info format"))
			c.Abort()
			return
		}

		// Construir contexto de permissão se necessário
		permContext := buildPermissionContext(c, userInfo)

		// Verificar permissão específica
		hasPermission, err := permissionService.HasPermission(c.Request.Context(), userInfo.ID, resource, action, permContext)
		if err != nil {
			slog.Error("Error checking specific permission",
				"userID", userInfo.ID, "resource", resource, "action", action, "error", err)
			httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Permission check failed"))
			c.Abort()
			return
		}

		if !hasPermission {
			slog.Warn("Specific permission denied",
				"userID", userInfo.ID, "resource", resource, "action", action)
			httperrors.SendHTTPErrorObj(c, coreutils.AuthorizationError("Insufficient permissions for this action"))
			c.Abort()
			return
		}

		slog.Debug("Specific permission granted",
			"userID", userInfo.ID, "resource", resource, "action", action)
		c.Next()
	})
}

// buildPermissionContext constrói o contexto de permissão baseado na requisição
func buildPermissionContext(c *gin.Context, userInfo usermodel.UserInfos) *permissionmodel.PermissionContext {
	context := permissionmodel.NewPermissionContext(userInfo.ID, userInfo.UserRoleID)

	// Adicionar metadados da requisição
	context.AddMetadata("request_ip", c.ClientIP()).
		AddMetadata("user_agent", c.GetHeader("User-Agent")).
		AddMetadata("method", c.Request.Method).
		AddMetadata("path", c.Request.URL.Path)

	// Adicionar parâmetros da URL ao contexto
	for _, param := range c.Params {
		context.AddMetadata("param_"+param.Key, param.Value)
	}

	// Adicionar query parameters relevantes
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			context.AddMetadata("query_"+key, values[0])
		}
	}

	return context
}

// isPermissionCheckRequired verifica se o endpoint precisa de verificação de permissão
func isPermissionCheckRequired(method, path string) bool {
	// Use a single source of truth for public endpoints
	// Method param kept for possible future expansion (e.g., method-specific rules)
	_ = method
	return !coreutils.IsPublicEndpoint(path)
}

// Helper functions for specific permission checks

// RequireListingPermission helper para permissões de listing
func RequireListingPermission(permissionService permissionservice.PermissionServiceInterface, action string) gin.HandlerFunc {
	return RequirePermission(permissionService, "listing", action)
}

// RequireUserPermission helper para permissões de usuário
func RequireUserPermission(permissionService permissionservice.PermissionServiceInterface, action string) gin.HandlerFunc {
	return RequirePermission(permissionService, "user", action)
}

// RequireComplexPermission helper para permissões de complexo
func RequireComplexPermission(permissionService permissionservice.PermissionServiceInterface, action string) gin.HandlerFunc {
	return RequirePermission(permissionService, "complex", action)
}

// RequireAdminPermission helper para permissões administrativas
func RequireAdminPermission(permissionService permissionservice.PermissionServiceInterface) gin.HandlerFunc {
	return RequirePermission(permissionService, "admin", "access")
}

// RequireOwnershipOrAdmin verifica se o usuário é dono do recurso ou admin
func RequireOwnershipOrAdmin(permissionService permissionservice.PermissionServiceInterface, resourceType string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		userInfoInterface, exists := c.Get("userInfo")
		if !exists {
			httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("User info not found"))
			c.Abort()
			return
		}

		userInfo, ok := userInfoInterface.(usermodel.UserInfos)
		if !ok {
			httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("Invalid user info format"))
			c.Abort()
			return
		}

		// Verificar se é admin primeiro
		permContext := buildPermissionContext(c, userInfo)
		hasAdminPermission, err := permissionService.HasPermission(c.Request.Context(), userInfo.ID, "admin", "access", permContext)
		if err == nil && hasAdminPermission {
			slog.Debug("Admin access granted", "userID", userInfo.ID, "resource", resourceType)
			c.Next()
			return
		}

		// Verificar ownership baseado no ID do recurso
		resourceIDStr := c.Param("id")
		if resourceIDStr == "" {
			httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("id", "Resource ID required"))
			c.Abort()
			return
		}

		resourceID, err := strconv.ParseInt(resourceIDStr, 10, 64)
		if err != nil {
			httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("id", "Invalid resource ID"))
			c.Abort()
			return
		}

		// Adicionar resource ID ao contexto para verificação de ownership
		permContext.AddMetadata("resource_id", resourceID).
			AddMetadata("resource_type", resourceType) // Verificar permissão de ownership
		hasOwnershipPermission, err := permissionService.HasPermission(c.Request.Context(), userInfo.ID, resourceType, "own", permContext)
		if err != nil {
			slog.Error("Error checking ownership permission",
				"userID", userInfo.ID, "resourceType", resourceType, "resourceID", resourceID, "error", err)
			httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Permission check failed"))
			c.Abort()
			return
		}

		if !hasOwnershipPermission {
			slog.Warn("Ownership permission denied",
				"userID", userInfo.ID, "resourceType", resourceType, "resourceID", resourceID)
			httperrors.SendHTTPErrorObj(c, coreutils.AuthorizationError("Access denied: insufficient permissions"))
			c.Abort()
			return
		}

		slog.Debug("Ownership permission granted",
			"userID", userInfo.ID, "resourceType", resourceType, "resourceID", resourceID)
		c.Next()
	})
}
