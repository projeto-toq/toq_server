package userservices

import (
	"context"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// PurgeStaleDeviceTokens deletes device_tokens whose updated_at is older than maxAge.
func (s *userService) PurgeStaleDeviceTokens(ctx context.Context, maxAge time.Duration, limit int) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, derrors.Infra("trace", err)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if limit <= 0 {
		limit = 500
	}

	cutoff := time.Now().Add(-maxAge)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("user.purge_device_tokens.tx_start_error", "err", txErr)
		return 0, derrors.Infra("start transaction", txErr)
	}

	var deleted int64
	defer func() {
		if err != nil {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("user.purge_device_tokens.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	deleted, err = s.repo.DeleteDeviceTokensOlderThan(ctx, tx, cutoff, limit)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.purge_device_tokens.delete_error", "err", err, "cutoff", cutoff, "limit", limit)
		return 0, derrors.Infra("delete stale device tokens", err)
	}

	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("user.purge_device_tokens.tx_commit_error", "err", cmErr)
		return 0, derrors.Infra("commit transaction", cmErr)
	}

	logger.Debug("user.purge_device_tokens.success", "deleted", deleted, "cutoff", cutoff, "limit", limit)
	return deleted, nil
}
