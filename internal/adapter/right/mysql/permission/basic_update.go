package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// Update executa uma query UPDATE e retorna o número de linhas afetadas
func (pa *PermissionAdapter) Update(ctx context.Context, tx *sql.Tx, query string, args ...any) (rowsAffected int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("mysqlpermissionadapter/Update: error preparing statement", "error", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		slog.Error("mysqlpermissionadapter/Update: error executing statement", "error", err)
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		slog.Error("mysqlpermissionadapter/Update: error getting rows affected", "error", err)
		return 0, err
	}

	return rowsAffected, nil
}
