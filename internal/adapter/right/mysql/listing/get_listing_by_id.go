package mysqllistingadapter

import (
	"context"
	"database/sql"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetListingByID(ctx context.Context, tx *sql.Tx, listingID int64) (listing listingmodel.ListingInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT *	FROM listings 
				WHERE id = ? AND deleted = 0;`

	return la.GetListingByQuery(ctx, tx, query, listingID)

}
