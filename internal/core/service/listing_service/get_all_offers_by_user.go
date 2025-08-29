package listingservices

import (
	"context"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ls *listingService) GetAllOffersByUser(ctx context.Context, userID int64) (offers []listingmodel.OfferInterface, err error) {
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
