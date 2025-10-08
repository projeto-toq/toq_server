package globalservice

import (
	"context"
	"database/sql"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (gs *globalService) StartTransaction(ctx context.Context) (tx *sql.Tx, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err = gs.globalRepo.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.tx_start_error", "err", err)
		return nil, utils.InternalError("")
	}

	return
}

// StartReadOnlyTransaction starts a read-only DB transaction via the global repository.
// Use this for pure read flows to minimize locking and overhead.
func (gs *globalService) StartReadOnlyTransaction(ctx context.Context) (tx *sql.Tx, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err = gs.globalRepo.StartReadOnlyTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.tx_start_readonly_error", "err", err)
		return nil, utils.InternalError("")
	}

	return
}

func (gs *globalService) RollbackTransaction(ctx context.Context, tx *sql.Tx) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	err = gs.globalRepo.RollbackTransaction(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.tx_rollback_error", "err", err)
		return utils.InternalError("")
	}

	return
}

func (gs *globalService) CommitTransaction(ctx context.Context, tx *sql.Tx) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	err = gs.globalRepo.CommitTransaction(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.tx_commit_error", "err", err)
		return utils.InternalError("")
	}

	return
}
