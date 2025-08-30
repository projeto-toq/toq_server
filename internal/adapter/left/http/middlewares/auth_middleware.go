package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	goroutines "github.com/giulio-alfieri/toq_server/internal/core/go_routines"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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
			c.JSON(http.StatusUnauthorized, utils.AuthenticationError("Authorization header required"))
			c.Abort()
			return
		}

		// Verify Bearer token format
		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) < 2 || tokenParts[1] == "" {
			c.JSON(http.StatusUnauthorized, utils.AuthenticationError("Invalid authorization token format"))
			c.Abort()
			return
		}

		token := tokenParts[1]
		userInfo, err := validateAccessToken(token)
		if err != nil {
			slog.Warn("Invalid access token", "error", err, "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, utils.AuthenticationError("Invalid access token"))
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

// setRootUserContext sets the root user context for public endpoints
func setRootUserContext(c *gin.Context) {
	infos := usermodel.UserInfos{
		ID:            0,
		Role:          permissionmodel.RoleSlugRoot, // Use RoleSlug instead of UserRole
		ProfileStatus: false,
	}

	// Set context values for compatibility
	ctx := context.WithValue(c.Request.Context(), globalmodel.TokenKey, infos)
	ctx = context.WithValue(ctx, globalmodel.UserAgentKey, c.GetHeader("User-Agent"))
	ctx = context.WithValue(ctx, globalmodel.ClientIPKey, c.ClientIP())

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

	c.Request = c.Request.WithContext(ctx)

	// Set Gin context values for easy access in handlers
	c.Set("userInfo", userInfo)
	c.Set("userAgent", c.GetHeader("User-Agent"))
	c.Set("clientIP", c.ClientIP())
}

// convertRoleNumberToSlug converts legacy role numbers to role slugs
func convertRoleNumberToSlug(roleNumber float64) permissionmodel.RoleSlug {
	switch int(roleNumber) {
	case 1:
		return permissionmodel.RoleSlugRoot
	case 2:
		return permissionmodel.RoleSlugOwner
	case 3:
		return permissionmodel.RoleSlugRealtor
	case 4:
		return permissionmodel.RoleSlugAgency
	default:
		return permissionmodel.RoleSlugOwner // Default fallback
	}
}

// validateAccessToken validates the JWT token and extracts user information
func validateAccessToken(tokenString string) (usermodel.UserInfos, error) {
	// This is a simplified version - we'll need to implement the actual JWT validation
	// For now, using the same logic from the gRPC auth interceptor

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// TODO: Get secret from environment configuration
		return []byte("Senh@123"), nil // This should come from config
	})

	if err != nil || !token.Valid {
		return usermodel.UserInfos{}, err
	}

	// Extract user information from claims
	userID, ok := (*claims)["user_id"].(float64)
	if !ok {
		return usermodel.UserInfos{}, jwt.NewValidationError("invalid user_id claim", jwt.ValidationErrorClaimsInvalid)
	}

	role, ok := (*claims)["role"].(float64)
	if !ok {
		return usermodel.UserInfos{}, jwt.NewValidationError("invalid role claim", jwt.ValidationErrorClaimsInvalid)
	}

	profileStatus, _ := (*claims)["profile_status"].(bool)

	return usermodel.UserInfos{
		ID:            int64(userID),
		Role:          convertRoleNumberToSlug(role),
		ProfileStatus: profileStatus,
	}, nil
}

// GetUserInfoFromContext extracts user info from Gin context
func GetUserInfoFromContext(c *gin.Context) (usermodel.UserInfos, bool) {
	userInfo, exists := c.Get("userInfo")
	if !exists {
		return usermodel.UserInfos{}, false
	}

	info, ok := userInfo.(usermodel.UserInfos)
	return info, ok
}
