package globalservice

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (gs *globalService) GetConfiguration(ctx context.Context) (configuration map[string]string, err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	tx, err := gs.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("global.get_configuration.tx_start_error", "err", err)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := gs.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("global.get_configuration.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	configuration, err = gs.globalRepo.GetConfiguration(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("global.get_configuration.query_error", "err", err)
		return nil, utils.InternalError("")
	}

	err = gs.CommitTransaction(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("global.get_configuration.tx_commit_error", "err", err)
		return nil, utils.InternalError("")
	}

	return

}
