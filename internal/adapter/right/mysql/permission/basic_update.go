package mysqlpermissionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// Update executa uma query UPDATE e retorna o número de linhas afetadas
func (pa *PermissionAdapter) Update(ctx context.Context, tx *sql.Tx, query string, args ...any) (rowsAffected int64, err error) {
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
		logger.Error("mysql.permission.update.prepare_error", "error", err)
		return 0, fmt.Errorf("prepare permission update statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.update.exec_error", "error", err)
		return 0, fmt.Errorf("exec permission update statement: %w", err)
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.permission.update.rows_affected_error", "error", err)
		return 0, fmt.Errorf("rows affected permission update: %w", err)
	}

	return rowsAffected, nil
}
