package photosessionservices

import (
	"context"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CleanOldAgendaEntries deletes agenda entries older than cutoff that are not referenced by bookings.
// The deletion executes in a transaction and emits metrics/logs for observability.
func (s *photoSessionService) CleanOldAgendaEntries(ctx context.Context, cutoff time.Time, limit int) (int64, error) {
	ctx, end, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return 0, utils.InternalError("")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if limit <= 0 || limit > 5000 {
		logger.Warn("photosession.cleaner.agenda.invalid_limit", "limit", limit)
		limit = 5000
	}

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photosession.cleaner.agenda.tx_start_error", "err", err)
		return 0, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("photosession.cleaner.agenda.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	deleted, err := s.repo.DeleteOldAgendaEntries(ctx, tx, cutoff, limit)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photosession.cleaner.agenda.delete_error", "err", err)
		return 0, utils.InternalError("")
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photosession.cleaner.agenda.tx_commit_error", "err", err)
		return 0, utils.InternalError("")
	}
	committed = true

	if deleted > 0 {
		metricPhotoSessionAgendaDeleted.Add(float64(deleted))
		logger.Info("photosession.cleaner.agenda.deleted", "count", deleted)
	}

	return deleted, nil
}
