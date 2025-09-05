package validationservice

import "context"

// Service defines the public contract for validation maintenance operations.
type Service interface {
	// CleanExpiredValidations deletes expired validation rows within a transaction boundary.
	CleanExpiredValidations(ctx context.Context, limit int) (int64, error)
}
