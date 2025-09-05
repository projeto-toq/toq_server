package sessionservice

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// CleanExpiredSessions deletes expired sessions using a transaction and records metrics.
func (s *service) CleanExpiredSessions(ctx context.Context, limit int) (int64, error) {
	// tracing span for public entrypoint
	ctx, end, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return 0, utils.InternalError("")
	}
	defer end()

	// Defensive limit to avoid huge deletes by mistake
	if limit <= 0 || limit > 5000 {
		limit = 5000
	}

	// Start transaction
	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("start_transaction_failed", "err", err)
		return 0, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("rollback_transaction_failed", "err", rbErr)
			}
		}
	}()

	// Do the deletion
	deleted, err := s.repo.DeleteExpiredSessions(ctx, tx, limit)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("delete_expired_sessions_failed", "err", err)
		return 0, utils.InternalError("")
	}

	// Commit
	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("commit_transaction_failed", "err", err)
		return 0, utils.InternalError("")
	}
	committed = true

	// Metrics best-effort
	metricSessionCleanerDeleted.Add(float64(deleted))
	slog.Info("clean_expired_sessions_success", "deleted", deleted)
	return deleted, nil
}

// compile-time check
var _ Service = (*service)(nil)
