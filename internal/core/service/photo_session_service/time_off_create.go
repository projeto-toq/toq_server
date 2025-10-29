package photosessionservices

import (
	"context"
	"database/sql"
	"strings"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateTimeOff registers a new time-off entry.
func (s *photoSessionService) CreateTimeOff(ctx context.Context, input TimeOffInput) (uint64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.time_off.create.tx_start_error", "err", err)
		return 0, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rollbackErr := s.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
				utils.SetSpanError(ctx, rollbackErr)
				logger.Error("photo_session.time_off.create.tx_rollback_error", "err", rollbackErr)
			}
		}
	}()

	id, err := s.createTimeOffInternal(ctx, tx, input)
	if err != nil {
		return 0, err
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.time_off.create.tx_commit_error", "err", err)
		return 0, utils.InternalError("")
	}
	committed = true

	return id, nil
}

// CreateTimeOffWithTx registers a new time-off entry using an existing transaction.
func (s *photoSessionService) CreateTimeOffWithTx(ctx context.Context, tx *sql.Tx, input TimeOffInput) (uint64, error) {
	if tx == nil {
		return 0, utils.InternalError("")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	return s.createTimeOffInternal(ctx, tx, input)
}

func (s *photoSessionService) createTimeOffInternal(ctx context.Context, tx *sql.Tx, input TimeOffInput) (uint64, error) {
	if err := validateTimeOffInput(input); err != nil {
		return 0, err
	}

	loc, tzErr := resolveLocation(input.Timezone)
	if tzErr != nil {
		return 0, tzErr
	}

	entry := photosessionmodel.NewAgendaEntry()
	entry.SetPhotographerUserID(input.PhotographerID)
	entry.SetEntryType(photosessionmodel.AgendaEntryTypeTimeOff)
	entry.SetSource(photosessionmodel.AgendaEntrySourceManual)
	entry.SetStartsAt(input.StartDate.UTC())
	entry.SetEndsAt(input.EndDate.UTC())
	entry.SetBlocking(true)
	entry.SetTimezone(loc.String())
	if input.Reason != nil {
		if reason := strings.TrimSpace(*input.Reason); reason != "" {
			entry.SetReason(reason)
		}
	}

	ids, err := s.repo.CreateEntries(ctx, tx, []photosessionmodel.AgendaEntryInterface{entry})
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("photo_session.time_off.create.repo_error", "photographer_id", input.PhotographerID, "err", err)
		return 0, utils.InternalError("")
	}
	if len(ids) == 0 {
		return 0, utils.InternalError("")
	}

	return ids[0], nil
}
