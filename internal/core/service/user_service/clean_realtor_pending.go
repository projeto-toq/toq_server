package userservices

import (
	"context"
	"database/sql"
	"errors"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) CleanRealtorPending(ctx context.Context, realtor usermodel.UserInterface) (err error) {
	ctx, spanEnd, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("user.clean_realtor_pending.tracer_error", "err", terr)
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	offers, err := us.listingService.GetAllOffersByUser(ctx, realtor.GetID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.clean_realtor_pending.get_offers_error", "realtor_id", realtor.GetID(), "err", err)
		return utils.InternalError("Failed to get offers by user")
	}
	for _, offer := range offers {
		if err = us.listingService.CancelOffer(ctx, offer.ID()); err != nil {
			return err
		}
	}

	visits, err := us.listingService.GetAllVisitsByUser(ctx, realtor.GetID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.clean_realtor_pending.get_visits_error", "realtor_id", realtor.GetID(), "err", err)
		return utils.InternalError("Failed to get visits by user")
	}
	for _, visit := range visits {
		if err = us.listingService.CancelVisit(ctx, visit.ID()); err != nil {
			return err
		}
	}

	err = us.DeleteAgencyOfRealtor(ctx, realtor.GetID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		// Erros inesperados: propagar; se for infra, será marcado a jusante
		return err
	}

	return nil
}
