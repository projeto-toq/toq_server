package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) DeleteComplex(ctx context.Context, tx *sql.Tx, id int64) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := "DELETE FROM complex WHERE id = ? LIMIT 1;"

	result, err := ca.ExecContext(ctx, tx, "delete", query, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.delete.exec_error", "error", err, "complex_id", id)
		return 0, fmt.Errorf("delete complex: %w", err)
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.delete.rows_affected_error", "error", err, "complex_id", id)
		return 0, fmt.Errorf("complex delete rows affected: %w", err)
	}

	return deleted, nil
}
