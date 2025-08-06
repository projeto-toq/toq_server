package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) Update(ctx context.Context, tx *sql.Tx, query string, args ...any) (affected int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("mysqluseradapter/Update: error preparing statement", "error", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		slog.Error("mysqluseradapter/Update: error executing statement", "error", err)
		return 0, err
	}

	affected, err = result.RowsAffected()
	if err != nil {
		slog.Error("mysqluseradapter/Update: error getting rows affected", "error", err)
		return 0, err
	}

	return affected, nil
}
