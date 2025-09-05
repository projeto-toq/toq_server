package http_errors

import (
	"net/http"

	"log/slog"
	"os"

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
	// If it's a DomainError, prefer preserving existing source if present.
	if derr, ok := err.(coreutils.DomainError); ok {
		debug := os.Getenv("TOQ_DEBUG_ERROR_TRACE") == "true"
		var out coreutils.DomainError
		// Se já possui origem (DomainErrorWithSource), não re-empacotar
		if _, hasSource := err.(coreutils.DomainErrorWithSource); hasSource {
			out = derr
			if debug {
				slog.Debug("http_errors: domain_error_with_source", "status", derr.Code())
			}
		} else {
			out = coreutils.WrapDomainErrorWithSource(derr)
			if debug {
				slog.Debug("http_errors: domain_error_wrapped", "status", out.Code())
			}
		}
		// Anexar erro e marcar span
		c.Error(out) //nolint: errcheck
		coreutils.SetSpanError(c.Request.Context(), out)
		payload := gin.H{"code": out.Code(), "message": out.Message()}
		if d := out.Details(); d != nil {
			payload["details"] = d
		}
		c.JSON(out.Code(), payload)
		return
	}
	derr := coreutils.InternalError("")
	if os.Getenv("TOQ_DEBUG_ERROR_TRACE") == "true" {
		slog.Debug("http_errors: generic_error_wrapped_internal", "status", derr.Code())
	}
	// Anexar erro ao contexto, marcar span e responder
	c.Error(derr) //nolint: errcheck
	coreutils.SetSpanError(c.Request.Context(), derr)
	c.JSON(derr.Code(), gin.H{"code": derr.Code(), "message": derr.Message()})
}
