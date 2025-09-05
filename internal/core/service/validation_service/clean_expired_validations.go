package validationservice

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// CleanExpiredValidations deletes expired validation rows within a transaction boundary
func (s *service) CleanExpiredValidations(ctx context.Context, limit int) (int64, error) {
	// Create tracing span for public entrypoint
	ctx, end, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return 0, utils.InternalError("")
	}
	defer end()

	// Defensive default if caller passes invalid limit
	if limit <= 0 {
		slog.Warn("validation.cleaner.invalid_limit", "limit", limit)
		limit = 500
	}

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("validation.cleaner.tx_start_error", "err", err)
		return 0, utils.InternalError("")
	}
	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("validation.cleaner.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	n, err := s.repo.DeleteExpiredValidations(ctx, tx, limit)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("validation.cleaner.delete_error", "err", err)
		return 0, utils.InternalError("")
	}
	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		slog.Error("validation.cleaner.tx_commit_error", "err", cmErr)
		return 0, utils.InternalError("")
	}
	committed = true

	if n > 0 {
		slog.Info("validation.cleaner.deleted", "count", n)
		metricValidationCleanerDeleted.Add(float64(n))
	}
	return n, nil
}
