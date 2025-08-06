package grpclistingport

import (
	"context"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
)

func (lr *ListingHandler) GuaranteesPBToModel(ctx context.Context, inputGuarantees []int64, listingID int64) (guarantees []listingmodel.GuaranteeInterface) {
	// ctx, spanEnd, err := utils.GenerateTracer(ctx)
	// if err != nil {
	// 	return
	// }
	// defer spanEnd()

	if len(inputGuarantees) == 0 {
		return
	}

	// guarantees = make([]listingmodel.GuaranteeInterface, len(inputGuarantees))
	for i, inputGuarantee := range inputGuarantees {
		guarantee := listingmodel.NewGuarantee()
		guarantee.SetListingID(listingID)
		guarantee.SetPriority(uint8(i + 1))
		guarantee.SetGuarantee(listingmodel.GuaranteeType(inputGuarantee))
		guarantees = append(guarantees, guarantee)
	}
	return
}
