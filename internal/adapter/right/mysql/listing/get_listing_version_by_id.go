package mysqllistingadapter

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

func (la *ListingAdapter) GetListingVersionByID(ctx context.Context, tx *sql.Tx, versionID int64) (listingmodel.ListingInterface, error) {
	query := `
		SELECT ` + listingSelectColumns + `
		FROM listing_versions lv
		JOIN listing_identities li ON lv.listing_identity_id = li.id
		WHERE lv.id = ? AND lv.deleted = 0
		LIMIT 1
	`
	return la.GetListingByQuery(ctx, tx, query, versionID)
}
