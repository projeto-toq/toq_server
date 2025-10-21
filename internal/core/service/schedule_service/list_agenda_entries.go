package scheduleservices

import (
	"context"
	"database/sql"
	"errors"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) ListAgendaEntries(ctx context.Context, filter schedulemodel.AgendaDetailFilter) (schedulemodel.AgendaDetailResult, error) {
	if filter.OwnerID <= 0 {
		return schedulemodel.AgendaDetailResult{}, utils.ValidationError("ownerId", "ownerId must be greater than zero")
	}
	if filter.ListingID <= 0 {
		return schedulemodel.AgendaDetailResult{}, utils.ValidationError("listingId", "listingId must be greater than zero")
	}
	if err := validateRange(filter.Range.From, filter.Range.To); err != nil {
		return schedulemodel.AgendaDetailResult{}, err
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return schedulemodel.AgendaDetailResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.list_agenda_entries.tx_start_error", "err", txErr)
		return schedulemodel.AgendaDetailResult{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("schedule.list_agenda_entries.tx_rollback_error", "err", rbErr)
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingID(ctx, tx, filter.ListingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schedulemodel.AgendaDetailResult{}, utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.list_agenda_entries.get_agenda_error", "listing_id", filter.ListingID, "err", err)
		return schedulemodel.AgendaDetailResult{}, utils.InternalError("")
	}

	if agenda.OwnerID() != filter.OwnerID {
		return schedulemodel.AgendaDetailResult{}, utils.AuthorizationError("Owner does not match listing agenda")
	}

	result, err := s.scheduleRepo.ListAgendaEntries(ctx, tx, filter)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.list_agenda_entries.repo_error", "listing_id", filter.ListingID, "err", err)
		return schedulemodel.AgendaDetailResult{}, utils.InternalError("")
	}

	return result, nil
}
