package derrors

import "net/http"

// Kind classifies an error so adapters can map it to transport concerns (e.g., HTTP status).
type Kind int

const (
	// Infrastructure/internal failures: DB, network, IO, unexpected states
	KindInfra Kind = iota
	// Authentication required/failed
	KindAuth
	// Authorization denied
	KindForbidden
	// Domain conflicts (resource state conflict)
	KindConflict
	// Domain: resource is gone/expired
	KindGone
	// Domain: semantically invalid payload/content
	KindUnprocessable
	// Domain: resource not found
	KindNotFound
	// Domain: request is invalid (format/body/query)
	KindBadRequest
	// Domain: validation error (fields)
	KindValidation
	// Throttling/locking
	KindLocked
	KindTooManyRequests
)

// HTTPStatus maps an error kind to an HTTP status code.
func HTTPStatus(k Kind) int {
	switch k {
	case KindAuth:
		return http.StatusUnauthorized
	case KindForbidden:
		return http.StatusForbidden
	case KindConflict:
		return http.StatusConflict
	case KindGone:
		return http.StatusGone
	case KindUnprocessable:
		return http.StatusUnprocessableEntity
	case KindNotFound:
		return http.StatusNotFound
	case KindBadRequest:
		return http.StatusBadRequest
	case KindValidation:
		return http.StatusBadRequest
	case KindLocked:
		return http.StatusLocked
	case KindTooManyRequests:
		return http.StatusTooManyRequests
	case KindInfra:
		fallthrough
	default:
		return http.StatusInternalServerError
	}
}
