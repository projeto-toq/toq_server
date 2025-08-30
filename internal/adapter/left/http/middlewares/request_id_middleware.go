package middlewares

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/google/uuid"
)

// RequestIDMiddleware generates and sets a unique request ID for each HTTP request
func RequestIDMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Generate a unique request ID
		requestID := uuid.New().String()

		// Add request ID to the standard context for service layer compatibility
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, globalmodel.RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		// Add request ID to Gin context for easy access in handlers and middlewares
		c.Set("request_id", requestID)

		// Add request ID to response headers for debugging and tracing
		c.Header("X-Request-ID", requestID)

		// Log the request ID for debugging
		slog.Debug("Request ID generated", "request_id", requestID, "path", c.Request.URL.Path)

		// Continue with the request
		c.Next()
	})
}
