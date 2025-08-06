package mysqllistingadapter

import (
	"context"
	"database/sql"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (la *ListingAdapter) UpdateFinancingBlockers(ctx context.Context, tx *sql.Tx, blockers []listingmodel.FinancingBlockerInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	//check if there is any data to update
	if len(blockers) == 0 {
		return
	}

	// Remove all blocker from listing
	err = la.DeleteListingFinancingBlockers(ctx, tx, blockers[0].ListingID())
	if err != nil {
		//check if the error is not found, because it's ok if there is no row to delete
		if status.Code(err) != codes.NotFound {
			return
		}
	}

	// Insert the new blokers
	for _, blocker := range blockers {
		err = la.CreateFinancingBlocker(ctx, tx, blocker)
		if err != nil {
			return
		}
	}

	return
}
