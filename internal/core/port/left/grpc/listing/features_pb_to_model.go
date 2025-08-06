package grpclistingport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
)

func (lr *ListingHandler) FeaturesPBToModel(ctx context.Context, inputFeatures []*pb.Feature, listingID int64) (features []listingmodel.FeatureInterface) {
	// ctx, spanEnd, err := utils.GenerateTracer(ctx)
	// if err != nil {
	// 	return
	// }
	// defer spanEnd()

	if len(inputFeatures) == 0 {
		return
	}

	// features = make([]listingmodel.FeatureInterface, len(inputFeatures))
	for _, inputFeature := range inputFeatures {
		feature := listingmodel.NewFeature()
		feature.SetListingID(listingID)
		feature.SetFeatureID(inputFeature.Id)
		feature.SetQuantity(uint8(inputFeature.Qty))
		features = append(features, feature)
	}
	return
}
