package mysqllistingadapter

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// GetDraftVersionByListingIdentityID retrieves the draft version of a listing identity
//
// This function queries for a non-active version with status=1 (draft) for the specified
// listing identity. Returns sql.ErrNoRows if no draft version exists.
//
// Business Rules:
//   - Only versions with status=1 (draft) are considered
//   - Only versions where lv.id != li.active_version_id are considered
//   - Only non-deleted versions (deleted=0)
//
// Parameters:
//   - ctx: Context for tracing, cancellation, and logging
//   - tx: Database transaction (can be nil for standalone queries)
//   - listingIdentityID: Unique identifier of the listing identity
//
// Returns:
//   - version: ListingInterface containing draft version data (FULLY ENRICHED)
//   - error: sql.ErrNoRows if no draft found, or other database errors
func (la *ListingAdapter) GetDraftVersionByListingIdentityID(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (listingmodel.ListingInterface, error) {
	// Query only draft versions (status=1) that are NOT active
	// Note: Uses GetListingByQuery which enriches all satellite tables
	query := `
		SELECT ` + listingSelectColumns + `
		FROM listing_versions lv
		JOIN listing_identities li ON lv.listing_identity_id = li.id
		WHERE lv.listing_identity_id = ? 
		  AND lv.status = 1
		  AND lv.id != COALESCE(li.active_version_id, 0)
		  AND lv.deleted = 0
		LIMIT 1
	`
	return la.GetListingByQuery(ctx, tx, query, listingIdentityID)
}
