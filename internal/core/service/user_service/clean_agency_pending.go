package userservices

import (
	"context"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (us *userService) CleanAgencyPending(ctx context.Context, agency usermodel.UserInterface) (err error) {
	ctx, spanEnd, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("user.clean_agency_pending.tracer_error", "err", terr)
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	realtors, err := us.GetRealtorsByAgency(ctx, agency.GetID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.clean_agency_pending.get_realtors_error", "agency_id", agency.GetID(), "err", err)
		return utils.InternalError("Failed to get realtors by agency")
	}

	for _, realtor := range realtors {
		err = us.DeleteRealtorOfAgency(ctx, agency.GetID(), realtor.GetID())
		if err != nil {
			// Erros aqui pertencem ao domínio/listing/permission; preferimos propagar sem marcar span como erro, a menos que seja infra em cascata
			return err
		}
	}

	return nil
}
