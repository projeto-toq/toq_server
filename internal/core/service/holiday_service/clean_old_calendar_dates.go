package holidayservices

import (
	"context"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CleanOldCalendarDates deletes non-recurrent calendar dates older than cutoff within a transactional boundary.
//
// Flow:
//  1. Start tracing/logging for observability.
//  2. Coerce limit to a safe value to avoid accidental full-table deletes.
//  3. Run deletion inside a transaction using the repository to remove rows older than cutoff.
//  4. Commit and emit metrics/logs with the deleted count.
//
// Parameters:
//   - ctx: base context (tracing/logging propagated).
//   - cutoff: absolute cutoff timestamp; dates strictly before this are deleted.
//   - limit: maximum rows to delete in this batch; coerced to sane bounds when invalid.
//
// Returns the number of deleted rows or an infrastructure error wrapped as internal.
func (s *holidayService) CleanOldCalendarDates(ctx context.Context, cutoff time.Time, limit int) (int64, error) {
	ctx, end, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return 0, utils.InternalError("")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if limit <= 0 || limit > 5000 {
		logger.Warn("holiday.cleaner.invalid_limit", "limit", limit)
		limit = 5000
	}

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.cleaner.tx_start_error", "err", err)
		return 0, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("holiday.cleaner.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	deleted, err := s.repo.DeleteOldCalendarDates(ctx, tx, cutoff, limit)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.cleaner.delete_error", "err", err)
		return 0, utils.InternalError("")
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("holiday.cleaner.tx_commit_error", "err", err)
		return 0, utils.InternalError("")
	}
	committed = true

	if deleted > 0 {
		metricHolidayCalendarDatesDeleted.Add(float64(deleted))
		logger.Info("holiday.cleaner.deleted", "count", deleted)
	}

	return deleted, nil
}
