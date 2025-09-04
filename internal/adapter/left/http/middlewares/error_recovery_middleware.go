package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// ErrorRecoveryMiddleware recovers from panics, converts them into a structured
// DomainErrorWithSource, attaches to context for logging, marks the span with error,
// and returns a standardized JSON response.
func ErrorRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Create an internal error with source information
				derr := coreutils.InternalError("")

				// Attach error to gin context for the logging middleware
				c.Error(derr) //nolint: errcheck

				// Mark span as error if tracing is active
				coreutils.SetSpanError(c.Request.Context(), derr)

				// Ensure we always return a JSON response
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    derr.Code(),
					"message": derr.Message(),
				})
			}
		}()

		c.Next()
	}
}
