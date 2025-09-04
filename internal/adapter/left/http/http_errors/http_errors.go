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
	// Criar erro com source e delegar para o caminho único (Obj)
	derr := coreutils.NewHTTPErrorWithSource(statusCode, message)
	SendHTTPErrorObj(c, derr)
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
		// Envolver em erro com source para logging de origem
		wrapped := coreutils.WrapDomainErrorWithSource(derr)
		// Anexar erro e marcar span
		c.Error(wrapped) //nolint: errcheck
		coreutils.SetSpanError(c.Request.Context(), wrapped)
		payload := gin.H{"code": wrapped.Code(), "message": wrapped.Message()}
		if d := wrapped.Details(); d != nil {
			payload["details"] = d
		}
		c.JSON(wrapped.Code(), payload)
		return
	}
	derr := coreutils.InternalError("")
	// Anexar erro ao contexto, marcar span e responder
	c.Error(derr) //nolint: errcheck
	coreutils.SetSpanError(c.Request.Context(), derr)
	c.JSON(derr.Code(), gin.H{"code": derr.Code(), "message": derr.Message()})
}
