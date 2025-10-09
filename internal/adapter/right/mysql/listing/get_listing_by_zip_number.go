package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetListingByZipNumber(ctx context.Context, tx *sql.Tx, zip string, number string) (listing listingmodel.ListingInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	query := `SELECT *	FROM listings 
				WHERE zip_code = ? AND number = ? AND deleted = 0;`

	listing, err = la.GetListingByQuery(ctx, tx, query, zip, number)
	if err != nil {
		return nil, fmt.Errorf("get listing by zip number: %w", err)
	}

	return listing, nil

}
