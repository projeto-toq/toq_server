package listingservices

import (
	"context"

"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ls *listingService) RejectOffer(ctx context.Context, offerID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// propertyTypes, err := ls.csi.GetOptions(ctx, zipCode, number)
	// if err != nil {
	// 	return
	// }

	// types = ls.DecodePropertyTypes(ctx, propertyTypes)

	return
}
