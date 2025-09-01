package http_errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	coreutils "github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// SendHTTPError (compat) – keeps current call sites working by accepting status and message directly.
// The error code string is ignored to keep payload flat and consistent (code/message/details).
func SendHTTPError(c *gin.Context, statusCode int, _ string, message string) {
	if c == nil {
		return
	}
	derr := coreutils.NewHTTPError(statusCode, message)
	payload := gin.H{"code": derr.Code(), "message": derr.Message()}
	if d := derr.Details(); d != nil {
		payload["details"] = d
	}
	c.JSON(derr.Code(), payload)
}

// SendHTTPErrorObj serializes a DomainError to HTTP. Accepts any DomainError implementation.
func SendHTTPErrorObj(c *gin.Context, err error) {
	if c == nil {
		return
	}
	if err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": "Internal server error"})
		return
	}
	// If it's a DomainError, pass through; otherwise wrap as InternalError
	if derr, ok := err.(coreutils.DomainError); ok {
		payload := gin.H{"code": derr.Code(), "message": derr.Message()}
		if d := derr.Details(); d != nil {
			payload["details"] = d
		}
		c.JSON(derr.Code(), payload)
		return
	}
	derr := coreutils.InternalError("")
	c.JSON(derr.Code(), gin.H{"code": derr.Code(), "message": derr.Message()})
}
