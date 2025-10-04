package listingservices

import (
	"context"

	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ls *listingService) GetAllListingsByUser(ctx context.Context, userID int64) (listings []listingmodel.ListingInterface, err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	// propertyTypes, err := ls.csi.GetOptions(ctx, zipCode, number)
	// if err != nil {
	// 	return
	// }

	// types = ls.DecodePropertyTypes(ctx, propertyTypes)

	return
}
