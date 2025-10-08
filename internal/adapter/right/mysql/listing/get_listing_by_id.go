package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetListingByID(ctx context.Context, tx *sql.Tx, listingID int64) (listing listingmodel.ListingInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	query := `SELECT *	FROM listings 
				WHERE id = ? AND deleted = 0;`

	listing, err = la.GetListingByQuery(ctx, tx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("get listing by id: %w", err)
	}

	return listing, nil

}
