package grpclistingport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
)

func (lr *ListingHandler) ExchangePlacesPBToModel(ctx context.Context, inputPlaces []*pb.Place, listingID int64) (places []listingmodel.ExchangePlaceInterface) {
	// ctx, spanEnd, err := utils.GenerateTracer(ctx)
	// if err != nil {
	// 	return
	// }
	// defer spanEnd()

	if len(inputPlaces) == 0 {
		return
	}

	// places = make([]listingmodel.ExchangePlaceInterface, len(inputPlaces))
	for _, inputPlace := range inputPlaces {
		place := listingmodel.NewExchangePlace()
		place.SetListingID(listingID)
		place.SetNeighborhood(inputPlace.Neighborhood)
		place.SetCity(inputPlace.City)
		place.SetState(inputPlace.State)
		places = append(places, place)
	}
	return
}
