package utils

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
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

// gRPC to HTTP status code mapping
var grpcToHTTPStatusMap = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           http.StatusRequestTimeout,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusBadRequest,
	codes.Aborted:            http.StatusConflict,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusInternalServerError,
	codes.Unauthenticated:    http.StatusUnauthorized,
}

// GRPCCodeToHTTPStatus converts gRPC status codes to HTTP status codes
func GRPCCodeToHTTPStatus(code codes.Code) int {
	if httpStatus, exists := grpcToHTTPStatusMap[code]; exists {
		return httpStatus
	}
	return http.StatusInternalServerError
}

// ConvertGRPCErrorToHTTP converts a gRPC error to an HTTP error
func ConvertGRPCErrorToHTTP(grpcCode codes.Code, message string, details ...interface{}) *HTTPError {
	httpCode := GRPCCodeToHTTPStatus(grpcCode)
	return NewHTTPError(httpCode, message, details...)
}

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
