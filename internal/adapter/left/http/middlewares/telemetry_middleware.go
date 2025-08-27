package middlewares

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// TelemetryMiddleware adds OpenTelemetry tracing to HTTP requests
func TelemetryMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		ctx := c.Request.Context()

		// Generate tracer for the request
		ctx, spanEnd, err := utils.GenerateTracer(ctx)
		if err != nil {
			slog.Warn("Failed to generate tracer", "error", err, "path", c.Request.URL.Path)
			c.Next()
			return
		}
		defer spanEnd()

		// Update request context with tracing context
		c.Request = c.Request.WithContext(ctx)

		// Continue with the request
		c.Next()
	})
}
