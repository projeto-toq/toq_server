package photosessionservices

import (
	"context"
	"database/sql"
	"errors"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteTimeOff removes an existing time-off entry.
func (s *photoSessionService) DeleteTimeOff(ctx context.Context, input DeleteTimeOffInput) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.time_off.delete.tx_start_error", "err", err)
		return utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rollbackErr := s.globalService.RollbackTransaction(ctx, tx); rollbackErr != nil {
				utils.SetSpanError(ctx, rollbackErr)
				logger.Error("photo_session.time_off.delete.tx_rollback_error", "err", rollbackErr)
			}
		}
	}()

	if err := s.deleteTimeOffInternal(ctx, tx, input); err != nil {
		return err
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.time_off.delete.tx_commit_error", "err", err)
		return utils.InternalError("")
	}
	committed = true
	return nil
}

// DeleteTimeOffWithTx removes an existing time-off entry inside a transaction.
func (s *photoSessionService) DeleteTimeOffWithTx(ctx context.Context, tx *sql.Tx, input DeleteTimeOffInput) error {
	if tx == nil {
		return utils.InternalError("")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	return s.deleteTimeOffInternal(ctx, tx, input)
}

func (s *photoSessionService) deleteTimeOffInternal(ctx context.Context, tx *sql.Tx, input DeleteTimeOffInput) error {
	logger := utils.LoggerFromContext(ctx)

	if tx == nil {
		return utils.InternalError("")
	}
	if input.TimeOffID == 0 {
		return utils.ValidationError("timeOffId", "timeOffId must be greater than zero")
	}

	txEntry, err := s.repo.GetEntryByIDForUpdate(ctx, tx, input.TimeOffID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("Time off")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.time_off.delete.get_error", "time_off_id", input.TimeOffID, "err", err)
		return utils.InternalError("")
	}

	if txEntry.EntryType() != photosessionmodel.AgendaEntryTypeTimeOff || txEntry.PhotographerUserID() != input.PhotographerID {
		return utils.NotFoundError("Time off")
	}

	if err := s.repo.DeleteEntryByID(ctx, tx, input.TimeOffID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("Time off")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("photo_session.time_off.delete.repo_error", "time_off_id", input.TimeOffID, "err", err)
		return utils.InternalError("")
	}

	return nil
}
