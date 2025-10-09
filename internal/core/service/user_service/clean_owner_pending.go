package userservices

import (
	"context"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (us *userService) CleanOwnerPending(ctx context.Context, owner usermodel.UserInterface) (err error) {
	ctx, spanEnd, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("user.clean_owner_pending.tracer_error", "err", terr)
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	offers, err := us.listingService.GetAllOffersByUser(ctx, owner.GetID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.clean_owner_pending.get_offers_error", "owner_id", owner.GetID(), "err", err)
		return utils.InternalError("Failed to get offers by user")
	}
	for _, offer := range offers {
		if err = us.listingService.RejectOffer(ctx, offer.ID()); err != nil {
			return err
		}
	}

	visits, err := us.listingService.GetAllVisitsByUser(ctx, owner.GetID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.clean_owner_pending.get_visits_error", "owner_id", owner.GetID(), "err", err)
		return utils.InternalError("Failed to get visits by user")
	}
	for _, visit := range visits {
		if err = us.listingService.RejectVisit(ctx, visit.ID()); err != nil {
			return err
		}
	}

	listings, err := us.listingService.GetAllListingsByUser(ctx, owner.GetID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.clean_owner_pending.get_listings_error", "owner_id", owner.GetID(), "err", err)
		return utils.InternalError("Failed to get listings by user")
	}
	for _, listing := range listings {
		if err = us.listingService.DeleteListing(ctx, listing.ID()); err != nil {
			return err
		}
	}

	return nil
}
