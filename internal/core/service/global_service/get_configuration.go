package globalservice

import (
	"context"

"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (gs *globalService) GetConfiguration(ctx context.Context) (configuration map[string]string, err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := gs.StartTransaction(ctx)
	if err != nil {
		return
	}

	configuration, err = gs.globalRepo.GetConfiguration(ctx, tx)
	if err != nil {
		gs.RollbackTransaction(ctx, tx)
		return
	}

	err = gs.CommitTransaction(ctx, tx)
	if err != nil {
		gs.RollbackTransaction(ctx, tx)
		return
	}

	return

}
