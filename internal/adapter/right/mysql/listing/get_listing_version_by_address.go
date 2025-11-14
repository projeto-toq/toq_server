package mysqllistingadapter

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// GetListingVersionByAddress retrieves an active listing version by zipCode and number.
// Returns sql.ErrNoRows if no active listing exists for the address.
// Active listings have status NOT IN (StatusDraft=1, StatusClosed=13, StatusExpired=15, StatusArchived=16).
func (la *ListingAdapter) GetListingVersionByAddress(ctx context.Context, tx *sql.Tx, zipCode, number string) (listingmodel.ListingInterface, error) {
	query := `
		SELECT ` + listingSelectColumns + `
		FROM listing_versions lv
		JOIN listing_identities li ON lv.listing_identity_id = li.id
		WHERE lv.zip_code = ? 
		  AND lv.number = ? 
		  AND lv.deleted = 0 
		  AND li.deleted = 0
		  AND lv.status NOT IN (1, 13, 15, 16)
		LIMIT 1
	`
	return la.GetListingByQuery(ctx, tx, query, zipCode, number)
}
