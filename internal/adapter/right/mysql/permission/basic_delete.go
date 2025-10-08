package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// Delete executa uma query DELETE e retorna o número de linhas afetadas
func (pa *PermissionAdapter) Delete(ctx context.Context, tx *sql.Tx, query string, args ...any) (rowsAffected int64, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.delete.prepare_error", "error", err)
		return 0, fmt.Errorf("prepare permission delete statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.delete.exec_error", "error", err)
		return 0, fmt.Errorf("exec permission delete statement: %w", err)
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.delete.rows_affected_error", "error", err)
		return 0, fmt.Errorf("rows affected permission delete: %w", err)
	}

	return rowsAffected, nil
}
