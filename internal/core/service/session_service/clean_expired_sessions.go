package sessionservice

import (
	"context"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CleanExpiredSessions deletes expired sessions using a transaction and records metrics.
func (s *service) CleanExpiredSessions(ctx context.Context, limit int) (int64, error) {
	return s.cleanExpiredSessionsInternal(ctx, nil, limit)
}

// CleanExpiredSessionsBefore deletes sessions whose effective expiry is before cutoff.
func (s *service) CleanExpiredSessionsBefore(ctx context.Context, cutoff time.Time, limit int) (int64, error) {
	return s.cleanExpiredSessionsInternal(ctx, &cutoff, limit)
}

func (s *service) cleanExpiredSessionsInternal(ctx context.Context, cutoff *time.Time, limit int) (int64, error) {
	ctx, end, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return 0, utils.InternalError("")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if limit <= 0 || limit > 5000 {
		limit = 5000
	}

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("start_transaction_failed", "err", err)
		return 0, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("rollback_transaction_failed", "err", rbErr)
			}
		}
	}()

	var deleted int64
	if cutoff == nil {
		deleted, err = s.repo.DeleteExpiredSessions(ctx, tx, limit)
	} else {
		deleted, err = s.repo.DeleteExpiredSessionsBefore(ctx, tx, *cutoff, limit)
	}
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("delete_expired_sessions_failed", "err", err)
		return 0, utils.InternalError("")
	}

	if err := s.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("commit_transaction_failed", "err", err)
		return 0, utils.InternalError("")
	}
	committed = true

	metricSessionCleanerDeleted.Add(float64(deleted))
	logger.Debug("clean_expired_sessions_success", "deleted", deleted)
	return deleted, nil
}
