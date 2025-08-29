package mysqllistingadapter

import (
	"context"
	"database/sql"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetListingByZipNumber(ctx context.Context, tx *sql.Tx, zip string, number string) (listing listingmodel.ListingInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT *	FROM listings 
				WHERE zip_code = ? AND number = ? AND deleted = 0;`

	return la.GetListingByQuery(ctx, tx, query, zip, number)

}
