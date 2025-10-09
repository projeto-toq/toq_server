package sessionmysqladapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (sa *SessionAdapter) StartTransaction(ctx context.Context) (tx *sql.Tx, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	tx, err = sa.db.DB.BeginTx(ctx, nil)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.session.transaction.begin_error", "error", err)
		return nil, fmt.Errorf("start transaction: %w", err)
	}
	return tx, nil
}

func (sa *SessionAdapter) RollbackTransaction(ctx context.Context, tx *sql.Tx) (err error) {
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
		logger.Error("mysql.session.transaction.rollback_error", "error", err)
		return fmt.Errorf("rollback transaction: %w", err)
	}
	return nil
}

func (sa *SessionAdapter) CommitTransaction(ctx context.Context, tx *sql.Tx) (err error) {
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
		logger.Error("mysql.session.transaction.commit_error", "error", err)
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
