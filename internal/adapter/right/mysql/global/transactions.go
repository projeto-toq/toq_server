package mysqlglobaladapter

import (
	"context"
	"database/sql"
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
		return nil, utils.ErrInternalServer
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
		return utils.ErrInternalServer
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
		return utils.ErrInternalServer
	}
	return nil
}
