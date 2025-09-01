package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTPError represents a structured HTTP error response
type HTTPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

// NewHTTPError creates a new HTTP error
func NewHTTPError(code int, message string, details ...interface{}) *HTTPError {
	httpErr := &HTTPError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		httpErr.Details = details[0]
	}
	return httpErr
}

// SendHTTPError sends an HTTP error response using Gin
func SendHTTPError(c *gin.Context, statusCode int, errorCode, message string) {
	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"code":    errorCode,
			"message": message,
		},
	})
}

// Predefined HTTP errors
var (
	ErrBadRequest          = &HTTPError{Code: http.StatusBadRequest, Message: "Bad Request"}
	ErrUnauthorized        = &HTTPError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden           = &HTTPError{Code: http.StatusForbidden, Message: "Forbidden"}
	ErrNotFound            = &HTTPError{Code: http.StatusNotFound, Message: "Not Found"}
	ErrMethodNotAllowed    = &HTTPError{Code: http.StatusMethodNotAllowed, Message: "Method Not Allowed"}
	ErrConflict            = &HTTPError{Code: http.StatusConflict, Message: "Conflict"}
	ErrUnprocessableEntity = &HTTPError{Code: http.StatusUnprocessableEntity, Message: "Unprocessable Entity"}
	ErrTooManyRequests     = &HTTPError{Code: http.StatusTooManyRequests, Message: "Too Many Requests"}
	ErrInternalServer      = &HTTPError{Code: http.StatusInternalServerError, Message: "Internal Server Error"}
	ErrBadGateway          = &HTTPError{Code: http.StatusBadGateway, Message: "Bad Gateway"}
	ErrServiceUnavailable  = &HTTPError{Code: http.StatusServiceUnavailable, Message: "Service Unavailable"}
	ErrGatewayTimeout      = &HTTPError{Code: http.StatusGatewayTimeout, Message: "Gateway Timeout"}
)

// ValidationError creates a structured validation error
func ValidationError(field, message string) *HTTPError {
	return NewHTTPError(http.StatusBadRequest, "Validation failed", map[string]string{
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
