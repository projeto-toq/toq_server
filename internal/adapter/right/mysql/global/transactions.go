package mysqlglobaladapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ga *GlobalAdapter) StartTransaction(ctx context.Context) (tx *sql.Tx, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()
	tx, err = ga.db.DB.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("Error starting transaction", "error", err)
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	return tx, nil
}

// StartReadOnlyTransaction starts a read-only transaction to optimize read flows.
func (ga *GlobalAdapter) StartReadOnlyTransaction(ctx context.Context) (tx *sql.Tx, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()
	tx, err = ga.db.DB.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		slog.Error("Error starting read-only transaction", "error", err)
		return nil, fmt.Errorf("begin tx readonly: %w", err)
	}
	return tx, nil
}

func (ga *GlobalAdapter) RollbackTransaction(ctx context.Context, tx *sql.Tx) (err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()
	err = tx.Rollback()
	if err != nil {
		slog.Error("Error rolling back transaction", "error", err)
		return fmt.Errorf("rollback tx: %w", err)
	}
	return nil
}

func (ga *GlobalAdapter) CommitTransaction(ctx context.Context, tx *sql.Tx) (err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()
	err = tx.Commit()
	if err != nil {
		slog.Error("Error committing transaction", "error", err)
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
