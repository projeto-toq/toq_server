package mysqluseradapter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) StartTransaction(ctx context.Context) (tx *sql.Tx, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()
	tx, err = ua.db.DB.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("mysqluseradapter/StartTransaction: error starting transaction", "error", err)
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	return tx, nil
}

func (ua *UserAdapter) RollbackTransaction(ctx context.Context, tx *sql.Tx) (err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()
	err = tx.Rollback()
	if err != nil {
		slog.Error("mysqluseradapter/RollbackTransaction: error rolling back transaction", "error", err)
		return fmt.Errorf("rollback tx: %w", err)
	}
	return nil
}

func (ua *UserAdapter) CommitTransaction(ctx context.Context, tx *sql.Tx) (err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()
	err = tx.Commit()
	if err != nil {
		slog.Error("mysqluseradapter/CommitTransaction: error committing transaction", "error", err)
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
