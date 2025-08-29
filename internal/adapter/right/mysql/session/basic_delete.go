package sessionmysqladapter

import (
	"context"
	"database/sql"
	"log/slog"

"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (sa *SessionAdapter) Delete(ctx context.Context, tx *sql.Tx, query string, args ...any) (deleted int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("sessionmysqladapter/Delete: error preparing statement", "error", err)
		return
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		slog.Error("sessionmysqladapter/Delete: error executing statement", "error", err)
		return
	}

	deleted, err = result.RowsAffected()
	if err != nil {
		slog.Error("sessionmysqladapter/Delete: error getting rows affected", "error", err)
		return
	}

	return
}
