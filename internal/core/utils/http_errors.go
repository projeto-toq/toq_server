package utils

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
)

// DomainErrorWithSource extends DomainError to expose error origin and stack.
// It allows adapters (e.g., HTTP middlewares) to log the real call site.
// Exported docs in English; internal comments in Portuguese quando necessário.
type DomainErrorWithSource interface {
	DomainError
	// Location returns function name, file path, and line number of the error origin.
	Location() (function string, file string, line int)
	// Stack returns a short stack trace captured at creation time (best-effort, may be empty).
	Stack() []string
}

// DomainError defines a core-agnostic error interface used across services and repositories
// It intentionally has no dependency on HTTP handlers or web frameworks.
type DomainError interface {
	error
	Code() int
	Message() string
	Details() any
}

// HTTPError is a concrete implementation of DomainError used by the core
// It represents a structured error with code/message/details.
type HTTPError struct {
	statusCode int    `json:"-"`
	msg        string `json:"-"`
	// details keeps optional error context (field, validation info, etc.)
	details any `json:"-"`
	// source info – não serializamos para resposta HTTP, apenas para logs
	function string   `json:"-"`
	file     string   `json:"-"`
	line     int      `json:"-"`
	stack    []string `json:"-"`
}

// Ensure *HTTPError implements DomainError
var _ DomainError = (*HTTPError)(nil)

// Ensure *HTTPError implements DomainErrorWithSource
var _ DomainErrorWithSource = (*HTTPError)(nil)

// Error implements the error interface
func (e *HTTPError) Error() string { return fmt.Sprintf("HTTP %d: %s", e.statusCode, e.msg) }

// Code returns the numeric error code (commonly an HTTP status code)
func (e *HTTPError) Code() int { return e.statusCode }

// Message returns the human-readable message
func (e *HTTPError) Message() string { return e.msg }

// Details returns optional error context
func (e *HTTPError) Details() any { return e.details }

// NewHTTPError creates a new structured error instance
func NewHTTPError(code int, message string, details ...any) *HTTPError {
	httpErr := &HTTPError{statusCode: code, msg: message}
	if len(details) > 0 {
		httpErr.details = details[0]
	}
	return httpErr
}

// errorStackDepth controls how many frames of stack are captured. Default 1.
var errorStackDepth = 1

// SetErrorStackDepth sets the number of stack frames to capture for DomainErrorWithSource.
// Values <= 0 disable stack capture (only Location will be set).
func SetErrorStackDepth(depth int) { //nolint: revive
	if depth <= 0 {
		errorStackDepth = 0
		return
	}
	errorStackDepth = depth
}

// NewHTTPErrorWithSource creates a new structured error instance capturing origin information.
func NewHTTPErrorWithSource(code int, message string, details ...any) *HTTPError { //nolint: revive
	// Capturar o callsite do chamador externo (service), não deste utilitário
	// skip segue convenção do runtime.Caller. Aqui usamos 2 para pular:
	// 0: newHTTPErrorWithSourceSkip, 1: NewHTTPErrorWithSource, 2: callsite do service
	return newHTTPErrorWithSourceSkip(code, message, 2, details...)
}

// newHTTPErrorWithSourceSkip cria erro estruturado capturando a origem com controle de profundidade.
// skip segue a convenção de runtime.Caller: 0 = esta função, 1 = quem chamou esta função, 2 = quem chamou o chamador, etc.
func newHTTPErrorWithSourceSkip(code int, message string, skip int, details ...any) *HTTPError { //nolint: revive
	he := NewHTTPError(code, message, details...)
	// Captura do local de criação baseado no skip informado
	if pc, file, line, ok := runtime.Caller(skip); ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			he.function = fn.Name()
		}
		he.file = file
		he.line = line
	}
	// Captura de uma stack curta (melhor esforço)
	if errorStackDepth > 0 {
		// pular este frame e o de runtime.Caller
		frames := make([]uintptr, errorStackDepth)
		n := runtime.Callers(skip+1, frames)
		he.stack = make([]string, 0, n)
		for i := 0; i < n; i++ {
			if f := runtime.FuncForPC(frames[i]); f != nil {
				file, line := f.FileLine(frames[i])
				he.stack = append(he.stack, fmt.Sprintf("%s (%s:%d)", f.Name(), filepath.Base(file), line))
			}
		}
	}
	return he
}

// Location returns the function, file and line where the error was created.
func (e *HTTPError) Location() (function string, file string, line int) { //nolint: revive
	return e.function, e.file, e.line
}

// Stack returns a short captured stack trace, if available.
func (e *HTTPError) Stack() []string { //nolint: revive
	return e.stack
}

// Predefined common errors (without coupling to handlers)
var (
	// Deprecated: prefer constructors that capture source.
	ErrBadRequest          = &HTTPError{statusCode: http.StatusBadRequest, msg: "Bad Request"}
	ErrUnauthorized        = &HTTPError{statusCode: http.StatusUnauthorized, msg: "Unauthorized"}
	ErrForbidden           = &HTTPError{statusCode: http.StatusForbidden, msg: "Forbidden"}
	ErrNotFound            = &HTTPError{statusCode: http.StatusNotFound, msg: "Not Found"}
	ErrMethodNotAllowed    = &HTTPError{statusCode: http.StatusMethodNotAllowed, msg: "Method Not Allowed"}
	ErrConflict            = &HTTPError{statusCode: http.StatusConflict, msg: "Conflict"}
	ErrUnprocessableEntity = &HTTPError{statusCode: http.StatusUnprocessableEntity, msg: "Unprocessable Entity"}
	ErrTooManyRequests     = &HTTPError{statusCode: http.StatusTooManyRequests, msg: "Too Many Requests"}
	ErrInternalServer      = &HTTPError{statusCode: http.StatusInternalServerError, msg: "Internal Server Error"}
	ErrBadGateway          = &HTTPError{statusCode: http.StatusBadGateway, msg: "Bad Gateway"}
	ErrServiceUnavailable  = &HTTPError{statusCode: http.StatusServiceUnavailable, msg: "Service Unavailable"}
	ErrGatewayTimeout      = &HTTPError{statusCode: http.StatusGatewayTimeout, msg: "Gateway Timeout"}
)

// Authentication/Refresh token specific errors (401 Unauthorized)
var (
	// Returned when a refresh token is invalid (malformed, bad signature, bad claims, or user mismatch)
	ErrInvalidRefreshToken = NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	// Returned when refresh token reuse is detected (old token being used after rotation)
	ErrRefreshTokenReuseDetected = NewHTTPError(http.StatusUnauthorized, "Refresh token reuse detected")
	// Returned when the refresh token chain reached absolute expiry
	ErrRefreshTokenExpired = NewHTTPError(http.StatusUnauthorized, "Refresh token expired")
	// Returned when the refresh token rotation limit was exceeded
	ErrRefreshRotationLimitExceeded = NewHTTPError(http.StatusUnauthorized, "Refresh token rotation limit exceeded")
	// Returned when a session cannot be found for the provided refresh token
	ErrRefreshSessionNotFound = NewHTTPError(http.StatusUnauthorized, "Session not found for refresh token")
)

// Email change domain-specific errors
var (
	// Returned when there is no pending email change for the user
	ErrEmailChangeNotPending = NewHTTPError(http.StatusConflict, "Email change not pending")
	// Returned when the provided email change code does not match the stored one
	ErrEmailChangeCodeInvalid = NewHTTPError(http.StatusUnprocessableEntity, "Invalid email change code")
	// Returned when the provided email change code is expired
	ErrEmailChangeCodeExpired = NewHTTPError(http.StatusGone, "Email change code expired")
	// Returned when attempting to set an email that is already used by another account
	ErrEmailAlreadyInUse = NewHTTPError(http.StatusConflict, "Email already in use")
	// Returned when the new email is the same as the current email
	ErrSameEmailAsCurrent = NewHTTPError(http.StatusConflict, "New email is the same as current")
)

// User/Role integrity errors
var (
	// Returned when a flow that requires an active role finds none
	ErrUserActiveRoleMissing = NewHTTPError(http.StatusConflict, "Active role missing for user")
)

// Phone change domain-specific errors
var (
	// Returned when there is no pending phone change for the user
	ErrPhoneChangeNotPending = NewHTTPError(http.StatusConflict, "Phone change not pending")
	// Returned when the provided phone change code does not match the stored one
	ErrPhoneChangeCodeInvalid = NewHTTPError(http.StatusUnprocessableEntity, "Invalid phone change code")
	// Returned when the provided phone change code is expired
	ErrPhoneChangeCodeExpired = NewHTTPError(http.StatusGone, "Phone change code expired")
	// Returned when attempting to set a phone that is already used by another account
	ErrPhoneAlreadyInUse = NewHTTPError(http.StatusConflict, "Phone already in use")
	// Returned when the new phone is the same as the current phone
	ErrSamePhoneAsCurrent = NewHTTPError(http.StatusConflict, "New phone is the same as current")
)

// Password change domain-specific errors
var (
	// Returned when there is no pending password change for the user
	ErrPasswordChangeNotPending = NewHTTPError(http.StatusConflict, "Password change not pending")
	// Returned when the provided password change code does not match the stored one
	ErrPasswordChangeCodeInvalid = NewHTTPError(http.StatusUnprocessableEntity, "Invalid password change code")
	// Returned when the provided password change code is expired
	ErrPasswordChangeCodeExpired = NewHTTPError(http.StatusGone, "Password change code expired")
)

// Photographer sessions specific errors
var (
	ErrPhotographerSlotUnavailable    = NewHTTPError(http.StatusConflict, "Photographer slot unavailable")
	ErrPhotographerReservationExpired = NewHTTPError(http.StatusGone, "Photographer slot reservation expired")
	ErrListingNotEligibleForPhoto     = NewHTTPError(http.StatusConflict, "Listing not eligible for photo session")
	ErrPhotoSessionNotCancelable      = NewHTTPError(http.StatusConflict, "Photo session cannot be cancelled")
)

// ValidationError creates a structured validation error
func ValidationError(field, message string) *HTTPError {
	if message == "" {
		message = "Validation failed"
	}
	return NewHTTPErrorWithSource(http.StatusBadRequest, message, map[string]string{
		"field":   field,
		"message": message,
	})
}

// AuthenticationError creates an authentication error
func AuthenticationError(message string) *HTTPError {
	if message == "" {
		message = "Authentication required"
	}
	return NewHTTPErrorWithSource(http.StatusUnauthorized, message)
}

// AuthorizationError creates an authorization error
func AuthorizationError(message string) *HTTPError {
	if message == "" {
		message = "Insufficient permissions"
	}
	return NewHTTPErrorWithSource(http.StatusForbidden, message)
}

// InternalError creates an internal server error
func InternalError(message string) *HTTPError {
	if message == "" {
		message = "Internal server error"
	}
	return NewHTTPErrorWithSource(http.StatusInternalServerError, message)
}

// NotFoundError creates a not found error
func NotFoundError(resource string) *HTTPError {
	message := "Resource not found"
	if resource != "" {
		message = fmt.Sprintf("%s not found", resource)
	}
	return NewHTTPErrorWithSource(http.StatusNotFound, message)
}

// ConflictError creates a conflict error
func ConflictError(message string) *HTTPError {
	if message == "" {
		message = "Resource conflict"
	}
	return NewHTTPErrorWithSource(http.StatusConflict, message)
}

// UserBlockedError creates a user blocked error (423 Locked)
func UserBlockedError(message string) *HTTPError {
	if message == "" {
		message = "Account temporarily blocked due to security measures"
	}
	return NewHTTPErrorWithSource(http.StatusLocked, message)
}

// TooManyAttemptsError creates a rate limiting error
func TooManyAttemptsError(message string) *HTTPError {
	if message == "" {
		message = "Too many failed attempts"
	}
	return NewHTTPErrorWithSource(http.StatusTooManyRequests, message)
}

// InvalidRequestFormatError creates a request format error
func InvalidRequestFormatError(message string) *HTTPError {
	if message == "" {
		message = "Invalid request format"
	}
	return NewHTTPErrorWithSource(http.StatusBadRequest, message)
}

// BadRequest creates a generic bad request error with source information.
func BadRequest(message string) *HTTPError { //nolint: revive
	if message == "" {
		message = "Bad Request"
	}
	return NewHTTPErrorWithSource(http.StatusBadRequest, message)
}

// WrapDomainErrorWithSource converts any DomainError to a new error instance that captures source.
// Útil nos adapters para preservar code/message/details e adicionar origem para logs.
func WrapDomainErrorWithSource(derr DomainError) *HTTPError { //nolint: revive
	if derr == nil {
		return newHTTPErrorWithSourceSkip(http.StatusInternalServerError, "Internal server error", 2)
	}
	// Copia detalhes se houver
	if d := derr.Details(); d != nil {
		return newHTTPErrorWithSourceSkip(derr.Code(), derr.Message(), 2, d)
	}
	return newHTTPErrorWithSourceSkip(derr.Code(), derr.Message(), 2)
}
