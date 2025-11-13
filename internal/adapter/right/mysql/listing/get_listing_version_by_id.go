package mysqllistingadapter

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

func (la *ListingAdapter) GetListingVersionByID(ctx context.Context, tx *sql.Tx, versionID int64) (listingmodel.ListingInterface, error) {
	return la.GetListingByID(ctx, tx, versionID)
}
