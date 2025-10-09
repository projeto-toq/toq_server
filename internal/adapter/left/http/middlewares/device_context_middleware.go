package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

// DeviceContextMiddleware extracts device id from headers and stores it in context.
// Header precedence: X-Device-Id, then optional fallback to empty if not provided.
func DeviceContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.GetHeader("X-Device-Id")

		// Expose in Gin context for convenience
		if deviceID != "" {
			c.Set("deviceID", deviceID)
		}

		// Propagate to request context for services
		ctx := context.WithValue(c.Request.Context(), globalmodel.DeviceIDKey, deviceID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
