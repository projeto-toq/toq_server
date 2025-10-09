package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
)

// RequireDeviceIDMiddleware ensures presence and validity of X-Device-Id header.
func RequireDeviceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.GetHeader("X-Device-Id")
		if deviceID == "" {
			httperrors.SendHTTPError(c, http.StatusBadRequest, "MISSING_DEVICE_ID", "X-Device-Id header is required")
			c.Abort()
			return
		}
		if _, err := uuid.Parse(deviceID); err != nil {
			httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_DEVICE_ID", "X-Device-Id must be a valid UUID")
			c.Abort()
			return
		}
		c.Next()
	}
}
