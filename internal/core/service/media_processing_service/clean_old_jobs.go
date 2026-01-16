package mediaprocessingservice

import (
	"context"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CleanOldJobs deletes terminal media processing jobs older than cutoff, capped by limit.
//
// Flow:
//  1. Initialize tracing/logging for observability.
//  2. Normalize limit to avoid unbounded deletions.
//  3. Run repository deletion inside a transaction.
//  4. Commit and emit metrics/logs with deleted count.
func (s *mediaProcessingService) CleanOldJobs(ctx context.Context, cutoff time.Time, limit int) (int64, error) {
	ctx, end, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return 0, utils.InternalError("")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if limit <= 0 || limit > 5000 {
		logger.Warn("media.cleaner.invalid_limit", "limit", limit)
		limit = 5000
	}

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("media.cleaner.tx_start_error", "err", err)
		return 0, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("media.cleaner.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	deleted, err := s.repo.DeleteOldJobs(ctx, tx, cutoff, limit)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("media.cleaner.delete_error", "err", err)
		return 0, utils.InternalError("")
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("media.cleaner.tx_commit_error", "err", err)
		return 0, utils.InternalError("")
	}
	committed = true

	if deleted > 0 {
		metricMediaJobsCleanerDeleted.Add(float64(deleted))
		logger.Info("media.cleaner.deleted", "count", deleted)
	}

	return deleted, nil
}
