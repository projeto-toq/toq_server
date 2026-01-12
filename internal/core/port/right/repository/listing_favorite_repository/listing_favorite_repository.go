package listingfavoriterepository

import (
	"context"
	"database/sql"
)

// FavoriteRepoPortInterface defines persistence operations for user ↔ listing favorites.
// Methods are idempotent where applicable to simplify upstream orchestration.
type FavoriteRepoPortInterface interface {
	// Add links a user to a listing identity as favorite. Must be idempotent.
	Add(ctx context.Context, tx *sql.Tx, userID, listingIdentityID int64) error

	// Remove unlinks a user from a listing identity. Idempotent: no error if absent.
	Remove(ctx context.Context, tx *sql.Tx, userID, listingIdentityID int64) error

	// ListByUser returns listing identity IDs favorited by the user plus total count (for pagination).
	ListByUser(ctx context.Context, tx *sql.Tx, userID int64, page, limit int) ([]int64, int64, error)

	// CountByListingIdentities returns favorites counts for the provided listing identity IDs.
	CountByListingIdentities(ctx context.Context, tx *sql.Tx, listingIdentityIDs []int64) (map[int64]int64, error)

	// GetUserFlags returns whether the given user has favorited each provided listing identity ID.
	GetUserFlags(ctx context.Context, tx *sql.Tx, listingIdentityIDs []int64, userID int64) (map[int64]bool, error)
}
