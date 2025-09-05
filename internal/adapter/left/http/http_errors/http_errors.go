package http_errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	derrors "github.com/giulio-alfieri/toq_server/internal/core/derrors"
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
	// 1) Novo fluxo: mapear derrors.KindError e sentinelas → HTTP
	status, message, details := mapErrorToHTTP(err)
	// Marcar span somente em 5xx para evitar ruído
	if status >= 500 {
		coreutils.SetSpanError(c.Request.Context(), err)
	}
	// Anexar erro ao contexto do gin para o middleware de logging
	c.Error(err) //nolint: errcheck
	c.JSON(status, gin.H{"code": status, "message": message, "details": details})
}

// mapErrorToHTTP converte erros do core (derrors) e legados para HTTP sem reempacotar.
func mapErrorToHTTP(err error) (status int, message string, details any) {
	// Novo: KindError com Kind → status
	if ke, ok := err.(derrors.KindError); ok {
		status = derrors.HTTPStatus(ke.Kind())
		message = ke.PublicMessage()
		details = ke.Details()
		if message == "" {
			message = http.StatusText(status)
		}
		return
	}
	// Sentinelas do domínio (telefone/email/role)
	switch {
	case errorsIs(err, derrors.ErrPhoneChangeNotPending):
		return http.StatusConflict, "Phone change not pending", nil
	case errorsIs(err, derrors.ErrPhoneChangeCodeInvalid):
		return http.StatusUnprocessableEntity, "Invalid phone change code", nil
	case errorsIs(err, derrors.ErrPhoneChangeCodeExpired):
		return http.StatusGone, "Phone change code expired", nil
	case errorsIs(err, derrors.ErrPhoneAlreadyInUse):
		return http.StatusConflict, "Phone already in use", nil
	case errorsIs(err, derrors.ErrUserActiveRoleMissing):
		return http.StatusConflict, "Active role missing for user", nil
	}
	// Legado: DomainError (utils) – preservar status/mensagem sem wrap
	if derr, ok := err.(coreutils.DomainError); ok {
		return derr.Code(), derr.Message(), derr.Details()
	}
	// Fallback: 500
	return http.StatusInternalServerError, "Internal server error", nil
}

// local helper to avoid importing "errors" at top due to existing imports order style
func errorsIs(err, target error) bool { //nolint: revive
	type is interface{ Is(error) bool }
	if x, ok := err.(is); ok {
		return x.Is(target)
	}
	// fallback
	for e := err; e != nil; {
		if e == target {
			return true
		}
		type unw interface{ Unwrap() error }
		if u, ok := e.(unw); ok {
			e = u.Unwrap()
			continue
		}
		break
	}
	return false
}
