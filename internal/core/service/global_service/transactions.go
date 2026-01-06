package globalservice

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// StartTransaction opens a read/write SQL transaction instrumented with tracing and structured logs.
// All infrastructure errors are converted to InternalError while being logged for observability.
func (gs *globalService) StartTransaction(ctx context.Context) (*sql.Tx, error) {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("global.transaction.tracer_error", "err", tracerErr)
		return nil, utils.InternalError("Failed to initialize transaction tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := gs.globalRepo.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.transaction.start_error", "err", err)
		return nil, utils.InternalError("Failed to start transaction")
	}

	return tx, nil
}

// StartReadOnlyTransaction starts a read-only DB transaction via the global repository.
// Use this for pure read flows to minimize locking and overhead.
func (gs *globalService) StartReadOnlyTransaction(ctx context.Context) (*sql.Tx, error) {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("global.transaction.ro.tracer_error", "err", tracerErr)
		return nil, utils.InternalError("Failed to initialize read-only transaction tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := gs.globalRepo.StartReadOnlyTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.transaction.readonly_start_error", "err", err)
		return nil, utils.InternalError("Failed to start read-only transaction")
	}

	return tx, nil
}

// RollbackTransaction ensures the provided transaction is rolled back with
// tracing/logging so infrastructure issues are observable by SRE dashboards.
func (gs *globalService) RollbackTransaction(ctx context.Context, tx *sql.Tx) error {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("global.transaction.rollback.tracer_error", "err", tracerErr)
		return utils.InternalError("Failed to initialize rollback tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err := gs.globalRepo.RollbackTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.transaction.rollback_error", "err", err)
		return utils.InternalError("Failed to rollback transaction")
	}

	return nil
}

// CommitTransaction commits the given transaction propagating tracing/logging metadata.
func (gs *globalService) CommitTransaction(ctx context.Context, tx *sql.Tx) error {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("global.transaction.commit.tracer_error", "err", tracerErr)
		return utils.InternalError("Failed to initialize commit tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err := gs.globalRepo.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("global.transaction.commit_error", "err", err)
		return utils.InternalError("Failed to commit transaction")
	}

	return nil
}
