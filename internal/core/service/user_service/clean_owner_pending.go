package userservices

import (
	"context"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func (us *userService) CleanOwnerPending(ctx context.Context, owner usermodel.UserInterface) (err error) {

	offers, err := us.listingService.GetAllOffersByUser(ctx, owner.GetID())
	if err != nil {
		return
	}
	for _, offer := range offers {
		err = us.listingService.RejectOffer(ctx, offer.ID())
		if err != nil {
			return
		}
	}

	visits, err := us.listingService.GetAllVisitsByUser(ctx, owner.GetID())
	if err != nil {
		return
	}
	for _, visit := range visits {
		err = us.listingService.RejectVisit(ctx, visit.ID())
		if err != nil {
			return
		}
	}

	listings, err := us.listingService.GetAllListingsByUser(ctx, owner.GetID())
	if err != nil {
		return
	}
	for _, listing := range listings {
		err = us.listingService.DeleteListing(ctx, listing.ID())
		if err != nil {
			return
		}
	}

	return
}
