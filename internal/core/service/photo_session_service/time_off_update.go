package photosessionservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateTimeOff updates a time-off entry and returns the mutated record.
func (s *photoSessionService) UpdateTimeOff(ctx context.Context, input UpdateTimeOffInput) (TimeOffDetailResult, error) {
	if input.TimeOffID == 0 {
		return TimeOffDetailResult{}, utils.ValidationError("timeOffId", "timeOffId must be greater than zero")
	}
	if err := validateTimeOffInput(TimeOffInput{
		PhotographerID:    input.PhotographerID,
		StartDate:         input.StartDate,
		EndDate:           input.EndDate,
		Reason:            input.Reason,
		Timezone:          input.Timezone,
		HolidayCalendarID: input.HolidayCalendarID,
		HorizonMonths:     input.HorizonMonths,
		WorkdayStartHour:  input.WorkdayStartHour,
		WorkdayEndHour:    input.WorkdayEndHour,
	}); err != nil {
		return TimeOffDetailResult{}, err
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
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

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_time_off.tx_start_error", "err", err)
		return TimeOffDetailResult{}, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rollbackErr := s.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
				utils.SetSpanError(ctx, rollbackErr)
				logger.Error("photo_session.update_time_off.tx_rollback_error", "err", rollbackErr)
			}
		}
	}()

	entry, err := s.repo.GetEntryByIDForUpdate(ctx, tx, input.TimeOffID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TimeOffDetailResult{}, utils.NotFoundError("Time off")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_time_off.get_error", "time_off_id", input.TimeOffID, "err", err)
		return TimeOffDetailResult{}, utils.InternalError("")
	}

	if entry.EntryType() != photosessionmodel.AgendaEntryTypeTimeOff || entry.PhotographerUserID() != input.PhotographerID {
		return TimeOffDetailResult{}, utils.NotFoundError("Time off")
	}

	entry.SetStartsAt(input.StartDate.UTC())
	entry.SetEndsAt(input.EndDate.UTC())
	entry.SetTimezone(loc.String())
	entry.SetSource(photosessionmodel.AgendaEntrySourceManual)
	entry.SetBlocking(true)
	entry.ClearReason()
	if input.Reason != nil {
		if reason := strings.TrimSpace(*input.Reason); reason != "" {
			entry.SetReason(reason)
		}
	}

	if err := s.repo.UpdateEntry(ctx, tx, entry); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_time_off.repo_error", "time_off_id", input.TimeOffID, "err", err)
		return TimeOffDetailResult{}, utils.InternalError("")
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.update_time_off.tx_commit_error", "err", err)
		return TimeOffDetailResult{}, utils.InternalError("")
	}
	committed = true

	entry.SetStartsAt(utils.ConvertToLocation(entry.StartsAt(), loc))
	entry.SetEndsAt(utils.ConvertToLocation(entry.EndsAt(), loc))

	return TimeOffDetailResult{TimeOff: entry, Timezone: loc.String()}, nil
}
