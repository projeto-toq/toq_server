package mysqllistingadapter

import (
	"context"
	"database/sql"
	"fmt"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (la *ListingAdapter) GetListingByID(ctx context.Context, tx *sql.Tx, listingID int64) (listing listingmodel.ListingInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	query := fmt.Sprintf(`SELECT
%s
FROM listing_versions lv
INNER JOIN listing_identities li ON li.id = lv.listing_identity_id
WHERE lv.id = ? AND lv.deleted = 0 AND li.deleted = 0;`, listingSelectColumns)

	listing, err = la.GetListingByQuery(ctx, tx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("get listing by id: %w", err)
	}

	return listing, nil

}
