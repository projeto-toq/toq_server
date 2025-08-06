package userservices

import (
	"context"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) CleanRealtorPending(ctx context.Context, realtor usermodel.UserInterface) (err error) {

	offers, err := us.listingService.GetAllOffersByUser(ctx, realtor.GetID())
	if err != nil {
		return
	}
	for _, offer := range offers {
		err = us.listingService.CancelOffer(ctx, offer.ID())
		if err != nil {
			return
		}
	}

	visits, err := us.listingService.GetAllVisitsByUser(ctx, realtor.GetID())
	if err != nil {
		return
	}
	for _, visit := range visits {
		err = us.listingService.CancelVisit(ctx, visit.ID())
		if err != nil {
			return
		}
	}

	err = us.DeleteAgencyOfRealtor(ctx, realtor.GetID())
	if err != nil && status.Code(err) != codes.NotFound {
		return
	}

	return nil
}
