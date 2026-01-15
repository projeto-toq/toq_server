package middlewares

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	goroutines "github.com/projeto-toq/toq_server/internal/core/go_routines"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	cacheport "github.com/projeto-toq/toq_server/internal/core/port/right/cache"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// AuthMiddleware handles JWT authentication with strict validation and blocklist enforcement.
func AuthMiddleware(activityTracker *goroutines.ActivityTracker, blocklist cacheport.TokenBlocklistPort) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		path := c.Request.URL.Path

		if !isAuthRequiredEndpoint(path) {
			setRootUserContext(c)
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("Authorization header required"))
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("auth", "missing_authorization")
			}
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) < 2 || tokenParts[1] == "" {
			httperrors.SendHTTPErrorObj(c, coreutils.AuthenticationError("Invalid authorization token format"))
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("auth", "invalid_format")
			}
			c.Abort()
			return
		}

		userInfo, jti, exp, err := validateAccessToken(c.Request.Context(), tokenParts[1], blocklist)
		if err != nil {
			httperrors.SendHTTPErrorObj(c, err)
			if mp := getMetricsAdapterFromGin(c); mp != nil {
				mp.IncrementErrors("auth", "invalid_token")
			}
			c.Abort()
			return
		}

		setUserContext(c, userInfo, jti, exp)

		if activityTracker != nil {
			activityTracker.TrackActivity(c.Request.Context(), userInfo.ID)
		}

		c.Next()
	})
}

func getMetricsAdapterFromGin(c *gin.Context) metricsport.MetricsPortInterface {
	if val, ok := c.Get("metricsAdapter"); ok {
		if mp, ok := val.(metricsport.MetricsPortInterface); ok {
			return mp
		}
	}
	return nil
}

func setRootUserContext(c *gin.Context) {
	infos := usermodel.UserInfos{
		ID:         0,
		UserRoleID: 0,
		RoleSlug:   permissionmodel.RoleSlugRoot,
	}

	ctx := context.WithValue(c.Request.Context(), globalmodel.TokenKey, infos)
	ctx = context.WithValue(ctx, globalmodel.UserAgentKey, c.GetHeader("User-Agent"))
	ctx = context.WithValue(ctx, globalmodel.ClientIPKey, c.ClientIP())
	ctx = coreutils.ContextWithLogger(ctx)

	c.Request = c.Request.WithContext(ctx)

	c.Set("userInfo", infos)
	c.Set("userAgent", c.GetHeader("User-Agent"))
	c.Set("clientIP", c.ClientIP())
}

func setUserContext(c *gin.Context, userInfo usermodel.UserInfos, tokenJTI string, expiresAt time.Time) {
	ctx := context.WithValue(c.Request.Context(), globalmodel.TokenKey, userInfo)
	ctx = context.WithValue(ctx, globalmodel.UserAgentKey, c.GetHeader("User-Agent"))
	ctx = context.WithValue(ctx, globalmodel.ClientIPKey, c.ClientIP())
	if tokenJTI != "" {
		ctx = context.WithValue(ctx, globalmodel.AccessTokenJTIKey, tokenJTI)
	}
	if !expiresAt.IsZero() {
		ctx = context.WithValue(ctx, globalmodel.AccessTokenExpiresAtKey, expiresAt)
	}
	ctx = coreutils.ContextWithLogger(ctx)

	c.Request = c.Request.WithContext(ctx)

	c.Set("userInfo", userInfo)
	c.Set("userAgent", c.GetHeader("User-Agent"))
	c.Set("clientIP", c.ClientIP())
	if tokenJTI != "" {
		c.Set("accessTokenJTI", tokenJTI)
	}
}

func validateAccessToken(ctx context.Context, tokenString string, blocklist cacheport.TokenBlocklistPort) (usermodel.UserInfos, string, time.Time, error) {
	claims := jwt.MapClaims{}
	parsed, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(globalmodel.GetJWTSecret()), nil
	})

	if err != nil || parsed == nil || !parsed.Valid {
		return usermodel.UserInfos{}, "", time.Time{}, coreutils.AuthenticationError("Invalid access token")
	}

	typ, _ := claims["typ"].(string)
	if typ != "access" {
		return usermodel.UserInfos{}, "", time.Time{}, coreutils.AuthenticationError("Invalid access token type")
	}

	jti, _ := claims["jti"].(string)
	if jti == "" {
		return usermodel.UserInfos{}, "", time.Time{}, coreutils.AuthenticationError("Invalid access token")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return usermodel.UserInfos{}, "", time.Time{}, coreutils.AuthenticationError("Invalid access token expiry")
	}
	expiresAt := time.Unix(int64(expFloat), 0)

	if blocklist != nil {
		blocked, blkErr := blocklist.Exists(ctx, jti)
		if blkErr != nil {
			return usermodel.UserInfos{}, "", time.Time{}, coreutils.InternalError("Failed to validate token")
		}
		if blocked {
			return usermodel.UserInfos{}, "", time.Time{}, coreutils.AuthenticationError("Access token revoked")
		}
	}

	userInfoClaim, ok := claims[string(globalmodel.TokenKey)]
	if !ok {
		return usermodel.UserInfos{}, "", time.Time{}, jwt.NewValidationError("missing user info claim", jwt.ValidationErrorClaimsInvalid)
	}

	userInfoMap, ok := userInfoClaim.(map[string]interface{})
	if !ok {
		return usermodel.UserInfos{}, "", time.Time{}, jwt.NewValidationError("invalid user info format", jwt.ValidationErrorClaimsInvalid)
	}

	userID, ok := userInfoMap["ID"].(float64)
	if !ok {
		return usermodel.UserInfos{}, "", time.Time{}, jwt.NewValidationError("invalid user ID claim", jwt.ValidationErrorClaimsInvalid)
	}

	userRoleID, ok := userInfoMap["UserRoleID"].(float64)
	if !ok {
		return usermodel.UserInfos{}, "", time.Time{}, jwt.NewValidationError("invalid user role ID claim", jwt.ValidationErrorClaimsInvalid)
	}

	var roleSlug permissionmodel.RoleSlug
	if rs, ok := userInfoMap["RoleSlug"].(string); ok {
		roleSlug = permissionmodel.RoleSlug(rs)
	}

	return usermodel.UserInfos{
		ID:         int64(userID),
		UserRoleID: int64(userRoleID),
		RoleSlug:   roleSlug,
	}, jti, expiresAt, nil
}

func GetUserInfoFromContext(c *gin.Context) (usermodel.UserInfos, bool) {
	userInfo, exists := c.Get("userInfo")
	if !exists {
		return usermodel.UserInfos{}, false
	}

	info, ok := userInfo.(usermodel.UserInfos)
	return info, ok
}
