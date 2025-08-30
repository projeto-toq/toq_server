package utils

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

// Context Utils para manipulação centralizada de contexto
// Segue melhores práticas Go e Google Style Guide

// GetUserInfoFromGinContext extrai informações do usuário do contexto Gin
// Retorna erro se usuário não estiver autenticado
func GetUserInfoFromGinContext(c *gin.Context) (usermodel.UserInfos, error) {
	userInfo, exists := c.Get("userInfo")
	if !exists {
		return usermodel.UserInfos{}, fmt.Errorf("user info not found in context")
	}

	info, ok := userInfo.(usermodel.UserInfos)
	if !ok {
		return usermodel.UserInfos{}, fmt.Errorf("invalid user info format in context")
	}

	return info, nil
}

// GetUserInfoFromContext extrai informações do usuário do contexto padrão
// Usado em services e outras camadas que recebem context.Context
func GetUserInfoFromContext(ctx context.Context) (usermodel.UserInfos, error) {
	userInfo := ctx.Value(globalmodel.TokenKey)
	if userInfo == nil {
		return usermodel.UserInfos{}, fmt.Errorf("user info not found in context")
	}

	info, ok := userInfo.(usermodel.UserInfos)
	if !ok {
		return usermodel.UserInfos{}, fmt.Errorf("invalid user info format in context")
	}

	return info, nil
}

// GetRequestIDFromContext extrai o Request ID do contexto
func GetRequestIDFromContext(ctx context.Context) string {
	requestID := ctx.Value(globalmodel.RequestIDKey)
	if requestID == nil {
		return ""
	}

	id, ok := requestID.(string)
	if !ok {
		return ""
	}

	return id
}

// GetRequestIDFromGinContext extrai o Request ID do contexto Gin
func GetRequestIDFromGinContext(c *gin.Context) string {
	requestID, exists := c.Get("request_id")
	if !exists {
		return ""
	}

	id, ok := requestID.(string)
	if !ok {
		return ""
	}

	return id
}

// GetClientIPFromGinContext extrai o IP do cliente do contexto Gin
func GetClientIPFromGinContext(c *gin.Context) string {
	clientIP, exists := c.Get("clientIP")
	if !exists {
		return c.ClientIP()
	}

	ip, ok := clientIP.(string)
	if !ok {
		return c.ClientIP()
	}

	return ip
}

// GetUserAgentFromGinContext extrai o User-Agent do contexto Gin
func GetUserAgentFromGinContext(c *gin.Context) string {
	userAgent, exists := c.Get("userAgent")
	if !exists {
		return c.GetHeader("User-Agent")
	}

	ua, ok := userAgent.(string)
	if !ok {
		return c.GetHeader("User-Agent")
	}

	return ua
}

// SetUserInContext adiciona informações do usuário ao contexto
func SetUserInContext(ctx context.Context, userInfo usermodel.UserInfos) context.Context {
	return context.WithValue(ctx, globalmodel.TokenKey, userInfo)
}

// SetRequestIDInContext adiciona Request ID ao contexto
func SetRequestIDInContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, globalmodel.RequestIDKey, requestID)
}

// EnrichContextWithRequestInfo enriquece o contexto com informações da requisição
func EnrichContextWithRequestInfo(ctx context.Context, c *gin.Context) context.Context {
	// Adicionar Request ID se disponível
	if requestID := GetRequestIDFromGinContext(c); requestID != "" {
		ctx = SetRequestIDInContext(ctx, requestID)
	}

	// Adicionar informações do usuário se disponível
	if userInfo, err := GetUserInfoFromGinContext(c); err == nil {
		ctx = SetUserInContext(ctx, userInfo)
	}

	// Adicionar metadados da requisição
	ctx = context.WithValue(ctx, globalmodel.UserAgentKey, GetUserAgentFromGinContext(c))
	ctx = context.WithValue(ctx, globalmodel.ClientIPKey, GetClientIPFromGinContext(c))

	return ctx
}

// RequireUserInContext valida que o usuário está autenticado no contexto
// Retorna erro se não encontrar usuário válido
func RequireUserInContext(ctx context.Context) (usermodel.UserInfos, error) {
	userInfo, err := GetUserInfoFromContext(ctx)
	if err != nil {
		return usermodel.UserInfos{}, fmt.Errorf("authentication required: %w", err)
	}

	if userInfo.ID == 0 {
		return usermodel.UserInfos{}, fmt.Errorf("invalid user ID in context")
	}

	return userInfo, nil
}

// IsAuthenticatedContext verifica se o contexto contém usuário autenticado
func IsAuthenticatedContext(ctx context.Context) bool {
	userInfo, err := GetUserInfoFromContext(ctx)
	if err != nil {
		return false
	}

	return userInfo.ID > 0
}

// IsPublicEndpoint verifica se um endpoint é público (não requer autenticação)
func IsPublicEndpoint(path string) bool {
	publicEndpoints := []string{
		"/api/v1/auth/owner",
		"/api/v1/auth/realtor",
		"/api/v1/auth/agency",
		"/api/v1/auth/signin",
		"/api/v1/auth/refresh",
		"/api/v1/auth/password/request",
		"/api/v1/auth/password/confirm",
		"/healthz",
		"/readyz",
		"/swagger/",
	}

	for _, endpoint := range publicEndpoints {
		if path == endpoint || (endpoint == "/swagger/" && len(path) > 8 && path[:9] == "/swagger/") {
			return true
		}
	}

	return false
}

// GetUserRoleFromContext extrai o role do usuário do contexto
func GetUserRoleFromContext(ctx context.Context) (usermodel.UserRole, error) {
	userInfo, err := GetUserInfoFromContext(ctx)
	if err != nil {
		return usermodel.UserRole(0), fmt.Errorf("failed to get user info: %w", err)
	}

	return userInfo.Role, nil
}

// HasRoleInContext verifica se o usuário tem um role específico
func HasRoleInContext(ctx context.Context, requiredRole usermodel.UserRole) bool {
	userRole, err := GetUserRoleFromContext(ctx)
	if err != nil {
		return false
	}

	return userRole == requiredRole
}

// IsAdminInContext verifica se o usuário é admin
func IsAdminInContext(ctx context.Context) bool {
	return HasRoleInContext(ctx, usermodel.RoleRoot)
}
