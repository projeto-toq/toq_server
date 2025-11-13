package scheduleservices

import (
	"context"
	"database/sql"
	"errors"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) GetAgendaByListingIdentityID(ctx context.Context, listingIdentityID int64) (schedulemodel.AgendaInterface, error) {
	if listingIdentityID <= 0 {
		return nil, utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.get_agenda.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("schedule.get_agenda.tx_rollback_error", "err", rbErr)
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, listingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.get_agenda.repo_error", "listing_identity_id", listingIdentityID, "err", err)
		return nil, utils.InternalError("")
	}

	return agenda, nil
}
