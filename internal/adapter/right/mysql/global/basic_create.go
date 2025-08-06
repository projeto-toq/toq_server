package mysqlglobaladapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ga *GlobalAdapter) Create(ctx context.Context, tx *sql.Tx, query string, args ...any) (id int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("mysqluseradapter/Create: error preparing statement", "error", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		slog.Error("mysqluseradapter/Create: error executing statement", "error", err)
		return id, nil
	}

	id, err = result.LastInsertId()
	if err != nil {
		slog.Error("mysqluseradapter/Create: error getting last insert ID", "error", err)
		return
	}

	return
}
