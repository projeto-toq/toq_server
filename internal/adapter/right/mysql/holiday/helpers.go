package mysqlholidayadapter

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

type sqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (a *HolidayAdapter) executor(tx *sql.Tx) sqlExecutor {
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

func defaultPagination(limit, page, max int) (int, int) {
	if limit <= 0 || limit > max {
		limit = max
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return limit, offset
}
