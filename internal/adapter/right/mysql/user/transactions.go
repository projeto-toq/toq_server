package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ua *UserAdapter) StartTransaction(ctx context.Context) (tx *sql.Tx, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	tx, err = ua.db.DB.BeginTx(ctx, nil)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.transaction.begin_error", "error", err)
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	return tx, nil
}

func (ua *UserAdapter) RollbackTransaction(ctx context.Context, tx *sql.Tx) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	err = tx.Rollback()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.transaction.rollback_error", "error", err)
		return fmt.Errorf("rollback tx: %w", err)
	}
	return nil
}

func (ua *UserAdapter) CommitTransaction(ctx context.Context, tx *sql.Tx) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	err = tx.Commit()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.user.transaction.commit_error", "error", err)
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
