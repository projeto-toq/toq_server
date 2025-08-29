package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"log/slog"

"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// Delete executa uma query DELETE e retorna o número de linhas afetadas
func (pa *PermissionAdapter) Delete(ctx context.Context, tx *sql.Tx, query string, args ...any) (rowsAffected int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("mysqlpermissionadapter/Delete: error preparing statement", "error", err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		slog.Error("mysqlpermissionadapter/Delete: error executing statement", "error", err)
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		slog.Error("mysqlpermissionadapter/Delete: error getting rows affected", "error", err)
		return 0, err
	}

	return rowsAffected, nil
}
