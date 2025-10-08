package middlewares

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	goroutines "github.com/giulio-alfieri/toq_server/internal/core/go_routines"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	metricsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/metrics"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/golang-jwt/jwt"
)

// AuthMiddleware handles JWT authentication
func AuthMiddleware(activityTracker *goroutines.ActivityTracker) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip authentication for public endpoints
		if !isAuthRequiredEndpoint(path) {
			setRootUserContext(c)
			c.Next()
			return
		}

		// Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("Authorization header required"))
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("auth", "missing_authorization")
			}
			c.Abort()
			return
		}

		// Verify Bearer token format
		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) < 2 || tokenParts[1] == "" {
			httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("Invalid authorization token format"))
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("auth", "invalid_format")
			}
			c.Abort()
			return
		}

		token := tokenParts[1]
		userInfo, err := validateAccessToken(token)
		if err != nil {
			httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("Invalid access token"))
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("auth", "invalid_token")
			}
			c.Abort()
			return
		}

		// Set user context
		setUserContext(c, userInfo)

		// Track user activity
		if activityTracker != nil {
			activityTracker.TrackActivity(c.Request.Context(), userInfo.ID)
		}

		c.Next()
	})
}

// Helper to get metrics adapter from Gin context (set by routes setup)
func getMetricsAdapterFromGin(c *gin.Context) metricsport.MetricsPortInterface {
	if val, ok := c.Get("metricsAdapter"); ok {
		if mp, ok := val.(metricsport.MetricsPortInterface); ok {
			return mp
		}
	}
	return nil
}

// setRootUserContext sets the root user context for public endpoints
func setRootUserContext(c *gin.Context) {
	infos := usermodel.UserInfos{
		ID:         0,
		UserRoleID: 0,
		RoleSlug:   permissionmodel.RoleSlugRoot,
	}

	// Set context values for compatibility
	ctx := context.WithValue(c.Request.Context(), globalmodel.TokenKey, infos)
	ctx = context.WithValue(ctx, globalmodel.UserAgentKey, c.GetHeader("User-Agent"))
	ctx = context.WithValue(ctx, globalmodel.ClientIPKey, c.ClientIP())
	ctx = coreutils.ContextWithLogger(ctx)

	c.Request = c.Request.WithContext(ctx)

	// Set Gin context values
	c.Set("userInfo", infos)
	c.Set("userAgent", c.GetHeader("User-Agent"))
	c.Set("clientIP", c.ClientIP())
}

// setUserContext sets the authenticated user context
func setUserContext(c *gin.Context, userInfo usermodel.UserInfos) {
	// Set context values for compatibility with existing service layer
	ctx := context.WithValue(c.Request.Context(), globalmodel.TokenKey, userInfo)
	ctx = context.WithValue(ctx, globalmodel.UserAgentKey, c.GetHeader("User-Agent"))
	ctx = context.WithValue(ctx, globalmodel.ClientIPKey, c.ClientIP())
	ctx = coreutils.ContextWithLogger(ctx)

	c.Request = c.Request.WithContext(ctx)

	// Set Gin context values for easy access in handlers
	c.Set("userInfo", userInfo)
	c.Set("userAgent", c.GetHeader("User-Agent"))
	c.Set("clientIP", c.ClientIP())
}

// validateAccessToken validates the JWT token and extracts user information
func validateAccessToken(tokenString string) (usermodel.UserInfos, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Get secret from global configuration
		return []byte(globalmodel.GetJWTSecret()), nil
	})

	if err != nil || !token.Valid {
		return usermodel.UserInfos{}, err
	}

	// Extract user information from new token structure
	userInfoClaim, ok := (*claims)[string(globalmodel.TokenKey)]
	if !ok {
		return usermodel.UserInfos{}, jwt.NewValidationError("missing user info claim", jwt.ValidationErrorClaimsInvalid)
	}

	// Parse UserInfos from claim
	userInfoMap, ok := userInfoClaim.(map[string]interface{})
	if !ok {
		return usermodel.UserInfos{}, jwt.NewValidationError("invalid user info format", jwt.ValidationErrorClaimsInvalid)
	}

	userID, ok := userInfoMap["ID"].(float64)
	if !ok {
		return usermodel.UserInfos{}, jwt.NewValidationError("invalid user ID claim", jwt.ValidationErrorClaimsInvalid)
	}

	userRoleID, ok := userInfoMap["UserRoleID"].(float64)
	if !ok {
		return usermodel.UserInfos{}, jwt.NewValidationError("invalid user role ID claim", jwt.ValidationErrorClaimsInvalid)
	}

	var roleSlug permissionmodel.RoleSlug
	if rs, ok := userInfoMap["RoleSlug"].(string); ok {
		roleSlug = permissionmodel.RoleSlug(rs)
	}

	return usermodel.UserInfos{
		ID:         int64(userID),
		UserRoleID: int64(userRoleID),
		RoleSlug:   roleSlug,
	}, nil
}

// Nota: Não logamos detalhes de tokens inválidos aqui. Métricas e o middleware de logging de request
// dão visibilidade adequada sem vazar metadados do JWT.

// GetUserInfoFromContext extracts user info from Gin context
func GetUserInfoFromContext(c *gin.Context) (usermodel.UserInfos, bool) {
	userInfo, exists := c.Get("userInfo")
	if !exists {
		return usermodel.UserInfos{}, false
	}

	info, ok := userInfo.(usermodel.UserInfos)
	return info, ok
}
