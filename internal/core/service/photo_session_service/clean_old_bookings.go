package photosessionservices

import (
	"context"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CleanOldBookings deletes photo session bookings in terminal statuses whose end time is before cutoff.
// It runs inside a transaction to ensure consistent deletion counts and emits metrics on success.
func (s *photoSessionService) CleanOldBookings(ctx context.Context, cutoff time.Time, limit int) (int64, error) {
	ctx, end, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return 0, utils.InternalError("")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if limit <= 0 || limit > 5000 {
		logger.Warn("photosession.cleaner.bookings.invalid_limit", "limit", limit)
		limit = 5000
	}

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photosession.cleaner.bookings.tx_start_error", "err", err)
		return 0, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("photosession.cleaner.bookings.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	deleted, err := s.repo.DeleteOldBookings(ctx, tx, cutoff, limit)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photosession.cleaner.bookings.delete_error", "err", err)
		return 0, utils.InternalError("")
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("photosession.cleaner.bookings.tx_commit_error", "err", err)
		return 0, utils.InternalError("")
	}
	committed = true

	if deleted > 0 {
		metricPhotoSessionBookingsDeleted.Add(float64(deleted))
		logger.Info("photosession.cleaner.bookings.deleted", "count", deleted)
	}

	return deleted, nil
}
