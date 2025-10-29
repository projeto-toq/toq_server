package photosessionservices

import (
	"context"
	"database/sql"
	"errors"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetTimeOffDetail fetches a specific time-off entry ensuring ownership.
func (s *photoSessionService) GetTimeOffDetail(ctx context.Context, input TimeOffDetailInput) (TimeOffDetailResult, error) {
	if input.TimeOffID == 0 {
		return TimeOffDetailResult{}, utils.ValidationError("timeOffId", "timeOffId must be greater than zero")
	}
	if input.PhotographerID == 0 {
		return TimeOffDetailResult{}, utils.ValidationError("photographerId", "photographerId must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "service.GetTimeOffDetail")
	if err != nil {
		return TimeOffDetailResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	loc, tzErr := resolveLocation(input.Timezone)
	if tzErr != nil {
		return TimeOffDetailResult{}, tzErr
	}

	tx, err := s.globalService.StartReadOnlyTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.get_time_off.tx_start_error", "err", err)
		return TimeOffDetailResult{}, utils.InternalError("")
	}
	defer func() {
		if rollbackErr := s.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
			utils.SetSpanError(ctx, rollbackErr)
			logger.Error("photo_session.get_time_off.tx_rollback_error", "err", rollbackErr)
		}
	}()

	entry, err := s.repo.GetEntryByID(ctx, tx, input.TimeOffID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TimeOffDetailResult{}, utils.NotFoundError("Time off")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.get_time_off.repo_error", "time_off_id", input.TimeOffID, "err", err)
		return TimeOffDetailResult{}, utils.InternalError("")
	}

	if entry.EntryType() != photosessionmodel.AgendaEntryTypeTimeOff || entry.PhotographerUserID() != input.PhotographerID {
		return TimeOffDetailResult{}, utils.NotFoundError("Time off")
	}

	entry.SetStartsAt(utils.ConvertToLocation(entry.StartsAt(), loc))
	entry.SetEndsAt(utils.ConvertToLocation(entry.EndsAt(), loc))

	return TimeOffDetailResult{TimeOff: entry, Timezone: loc.String()}, nil
}
