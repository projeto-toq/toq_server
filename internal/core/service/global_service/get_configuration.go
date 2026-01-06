package globalservice

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetConfiguration loads configuration key-value pairs using a read-only transaction.
// It guarantees tracing/logging coverage for each infrastructure interaction.
func (gs *globalService) GetConfiguration(ctx context.Context) (configuration map[string]string, err error) {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("global.get_configuration.tracer_error", "err", tracerErr)
		return nil, utils.InternalError("Failed to initialize configuration tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := gs.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("global.get_configuration.tx_start_error", "err", txErr)
		return nil, txErr
	}

	defer func() {
		if err != nil {
			if rbErr := gs.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("global.get_configuration.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	configuration, err = gs.globalRepo.GetConfiguration(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.get_configuration.query_error", "err", err)
		err = utils.InternalError("Failed to load configuration entries")
		return nil, err
	}

	if commitErr := gs.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("global.get_configuration.tx_commit_error", "err", commitErr)
		return nil, commitErr
	}

	return configuration, nil
}
