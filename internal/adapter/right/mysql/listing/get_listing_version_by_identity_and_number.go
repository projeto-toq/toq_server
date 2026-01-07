package mysqllistingadapter

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// GetListingVersionByIdentityAndNumber retrieves a listing version by identity ID and version number.
//
// Parameters:
//   - listingIdentityID: logical listing identifier (listing_identities.id)
//   - version: version number stored in listing_versions.version
//
// Returns sql.ErrNoRows when the version does not exist or is marked as deleted.
func (la *ListingAdapter) GetListingVersionByIdentityAndNumber(ctx context.Context, tx *sql.Tx, listingIdentityID int64, version uint8) (listingmodel.ListingInterface, error) {
	query := `
        SELECT ` + listingSelectColumns + `
        FROM listing_versions lv
        JOIN listing_identities li ON lv.listing_identity_id = li.id
        WHERE lv.listing_identity_id = ? AND lv.version = ? AND lv.deleted = 0
        LIMIT 1`

	return la.GetListingByQuery(ctx, tx, query, listingIdentityID, version)
}
