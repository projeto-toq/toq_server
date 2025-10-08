package globalservice

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (gs *globalService) GetConfiguration(ctx context.Context) (configuration map[string]string, err error) {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("global.get_configuration.tracer_error", "err", tracerErr)
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := gs.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.get_configuration.tx_start_error", "err", err)
		return nil, utils.InternalError("")
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
		return nil, utils.InternalError("")
	}

	err = gs.CommitTransaction(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.get_configuration.tx_commit_error", "err", err)
		return nil, utils.InternalError("")
	}

	return

}
