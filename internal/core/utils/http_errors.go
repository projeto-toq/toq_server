package utils

import (
	"fmt"
	"net/http"
)

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
}

// Ensure *HTTPError implements DomainError
var _ DomainError = (*HTTPError)(nil)

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

// Predefined common errors (without coupling to handlers)
var (
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

// ValidationError creates a structured validation error
func ValidationError(field, message string) *HTTPError {
	if message == "" {
		message = "Validation failed"
	}
	return NewHTTPError(http.StatusBadRequest, message, map[string]string{
		"field":   field,
		"message": message,
	})
}

// AuthenticationError creates an authentication error
func AuthenticationError(message string) *HTTPError {
	if message == "" {
		message = "Authentication required"
	}
	return NewHTTPError(http.StatusUnauthorized, message)
}

// AuthorizationError creates an authorization error
func AuthorizationError(message string) *HTTPError {
	if message == "" {
		message = "Insufficient permissions"
	}
	return NewHTTPError(http.StatusForbidden, message)
}

// InternalError creates an internal server error
func InternalError(message string) *HTTPError {
	if message == "" {
		message = "Internal server error"
	}
	return NewHTTPError(http.StatusInternalServerError, message)
}

// NotFoundError creates a not found error
func NotFoundError(resource string) *HTTPError {
	message := "Resource not found"
	if resource != "" {
		message = fmt.Sprintf("%s not found", resource)
	}
	return NewHTTPError(http.StatusNotFound, message)
}

// ConflictError creates a conflict error
func ConflictError(message string) *HTTPError {
	if message == "" {
		message = "Resource conflict"
	}
	return NewHTTPError(http.StatusConflict, message)
}

// UserBlockedError creates a user blocked error (423 Locked)
func UserBlockedError(message string) *HTTPError {
	if message == "" {
		message = "Account temporarily blocked due to security measures"
	}
	return NewHTTPError(http.StatusLocked, message)
}

// TooManyAttemptsError creates a rate limiting error
func TooManyAttemptsError(message string) *HTTPError {
	if message == "" {
		message = "Too many failed attempts"
	}
	return NewHTTPError(http.StatusTooManyRequests, message)
}

// InvalidRequestFormatError creates a request format error
func InvalidRequestFormatError(message string) *HTTPError {
	if message == "" {
		message = "Invalid request format"
	}
	return NewHTTPError(http.StatusBadRequest, message)
}
