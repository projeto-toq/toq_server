package listingviewrepository

import (
	"context"
	"database/sql"
)

// Repository defines persistence operations for listing view metrics.
// This port isolates view tracking storage so adapters can implement
// atomic counters without leaking database concerns to the service layer.
// All methods must be safe to call concurrently and should rely on
// transactional semantics to guarantee consistency.
type Repository interface {
	// IncrementAndGet atomically increments the view counter for a listing identity.
	// Returns the updated total after the increment.
	// Implementations must be idempotent per call and handle concurrent increments safely.
	IncrementAndGet(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (int64, error)

	// GetCount returns the current view counter for a listing identity.
	// Should return 0 when no record exists instead of sql.ErrNoRows.
	GetCount(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (int64, error)
}
