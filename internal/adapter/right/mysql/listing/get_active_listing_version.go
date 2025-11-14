package mysqllistingadapter

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// GetActiveListingVersion retrieves the currently active version for a listing identity.
// Returns sql.ErrNoRows if the identity doesn't exist or has no active version.
func (la *ListingAdapter) GetActiveListingVersion(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (listingmodel.ListingInterface, error) {
	query := `
		SELECT ` + listingSelectColumns + `
		FROM listing_versions lv
		JOIN listing_identities li ON lv.id = li.active_version_id
		WHERE li.id = ? 
		  AND li.deleted = 0 
		  AND lv.deleted = 0
		LIMIT 1
	`
	return la.GetListingByQuery(ctx, tx, query, listingIdentityID)
}
