package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) UpdateComplexSize(ctx context.Context, tx *sql.Tx, size complexmodel.ComplexSizeInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE complex_sizes SET
		size = ?,
		description = ?
	WHERE id = ?
	AND complex_id = ?
	LIMIT 1;`

	result, err := ca.ExecContext(
		ctx,
		tx,
		"update",
		query,
		size.Size(),
		nullableStringValue(size.Description()),
		size.ID(),
		size.ComplexID(),
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.size.update.exec_error", "error", err, "size_id", size.ID(), "complex_id", size.ComplexID())
		return 0, fmt.Errorf("update complex size: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.size.update.rows_affected_error", "error", err, "size_id", size.ID(), "complex_id", size.ComplexID())
		return 0, fmt.Errorf("complex size rows affected: %w", err)
	}

	return affected, nil
}
