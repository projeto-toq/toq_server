package middlewares

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	metricsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/metrics"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// TelemetryMiddleware adds OpenTelemetry tracing to HTTP requests
// Now also integrates with metrics collection
func TelemetryMiddleware(metricsAdapter metricsport.MetricsPortInterface) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		ctx := c.Request.Context()

		// Generate tracer for the request
		ctx, spanEnd, err := utils.GenerateTracer(ctx)
		if err != nil {
			slog.Warn("Failed to generate tracer", "error", err, "path", c.Request.URL.Path)
			// Continue even if tracing fails
		} else {
			// Update request context with tracing context
			c.Request = c.Request.WithContext(ctx)
			defer spanEnd()
		}

		// Apply metrics middleware if adapter is provided
		if metricsAdapter != nil {
			MetricsMiddleware(metricsAdapter)(c)
		} else {
			c.Next()
		}
	})
}
