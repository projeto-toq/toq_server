package globalservice

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (gs *globalService) GetConfiguration(ctx context.Context) (configuration map[string]string, err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	tx, err := gs.StartTransaction(ctx)
	if err != nil {
		slog.Error("global.get_configuration.tx_start_error", "err", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			if rbErr := gs.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("global.get_configuration.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	configuration, err = gs.globalRepo.GetConfiguration(ctx, tx)
	if err != nil {
		return nil, err
	}

	err = gs.CommitTransaction(ctx, tx)
	if err != nil {
		slog.Error("global.get_configuration.tx_commit_error", "err", err)
		return nil, err
	}

	return

}
