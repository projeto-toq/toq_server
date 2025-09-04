package utils

import (
	"context"
	"database/sql"
)

// MapRepositoryError converts low-level repository errors into DomainError.
// Use this in services to keep a consistent domain boundary without HTTP coupling.
// English docs for exported functions; Portuguese comments for internal notes.
func MapRepositoryError(err error, notFoundMessage string) *HTTPError {
	if err == nil {
		return nil
	}
	// Casos mais comuns e portáveis:
	if err == sql.ErrNoRows {
		if notFoundMessage == "" {
			return NotFoundError("")
		}
		return NotFoundError(notFoundMessage)
	}
	// Fallback seguro: InternalError genérico
	return InternalError("")
}

// MarkSpanIfError is a tiny helper to mark the current span when returning errors in services.
// It's safe to call with nil err.
func MarkSpanIfError(ctx context.Context, err error) {
	if err != nil {
		SetSpanError(ctx, err)
	}
}
