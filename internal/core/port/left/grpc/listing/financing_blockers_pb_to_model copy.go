package grpclistingport

import (
	"context"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
)

func (lr *ListingHandler) FinancingBlockersPBToModel(ctx context.Context, inputBlockers []int64, listingID int64) (blockers []listingmodel.FinancingBlockerInterface) {
	// ctx, spanEnd, err := utils.GenerateTracer(ctx)
	// if err != nil {
	// 	return
	// }
	// defer spanEnd()

	if len(inputBlockers) == 0 {
		return
	}

	// blockers = make([]listingmodel.FinancingBlockerInterface, len(inputBlockers))
	for _, inputBlocker := range inputBlockers {
		blocker := listingmodel.NewFinancingBlocker()
		blocker.SetListingID(listingID)
		blocker.SetBlocker(listingmodel.FinancingBlocker(inputBlocker))
		blockers = append(blockers, blocker)
	}
	return
}
