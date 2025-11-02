package mysqlphotosessionadapter

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

type sqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

func (a *PhotoSessionAdapter) executor(tx *sql.Tx) sqlExecutor {
	if tx != nil {
		return tx
	}
	return a.DB().GetDB()
}

func withTracer(ctx context.Context) (context.Context, func(), error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, spanEnd, nil
}
