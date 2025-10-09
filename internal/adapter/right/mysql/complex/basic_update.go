package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) Update(ctx context.Context, tx *sql.Tx, query string, args ...any) (affected int64, err error) {
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
		logger.Error("mysql.complex.update.prepare_error", "error", err)
		return 0, fmt.Errorf("prepare complex update statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.update.exec_error", "error", err)
		return 0, fmt.Errorf("exec complex update statement: %w", err)
	}

	affected, err = result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.update.rows_affected_error", "error", err)
		return 0, fmt.Errorf("rows affected complex update: %w", err)
	}

	return affected, nil
}
