package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// AccessControlMiddleware checks if the user has permission to access the endpoint
func AccessControlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip access control for public endpoints (signin, signup, etc.)
		if !isAccessControlRequired(c.Request.Method, c.Request.URL.Path) {
			c.Next()
			return
		}

		// Get user info from context (set by authentication middleware)
		userInfoInterface, exists := c.Get("userInfo")
		if !exists {
			utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "No user info found")
			c.Abort()
			return
		}

		userInfo, ok := userInfoInterface.(usermodel.UserInfos)
		if !ok {
			utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid user info")
			c.Abort()
			return
		}

		// Get user role from user info
		userRole := userInfo.Role

		// Get request method and path
		method := c.Request.Method
		path := c.Request.URL.Path

		// Check if user has access to this endpoint
		if !usermodel.IsHTTPEndpointAllowed(userRole, method, path) {
			utils.SendHTTPError(c, http.StatusForbidden, "FORBIDDEN",
				"User does not have permission to access this endpoint")
			c.Abort()
			return
		}

		// Continue to next handler
		c.Next()
	}
}
