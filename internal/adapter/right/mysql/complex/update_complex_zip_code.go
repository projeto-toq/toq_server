package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) UpdateComplexZipCode(ctx context.Context, tx *sql.Tx, zip complexmodel.ComplexZipCodeInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE complex_zip_codes SET
		zip_code = ?
	WHERE id = ?
	AND complex_id = ?
	LIMIT 1;`

	result, err := ca.ExecContext(
		ctx,
		tx,
		"update",
		query,
		zip.ZipCode(),
		zip.ID(),
		zip.ComplexID(),
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.zip_code.update.exec_error", "error", err, "zip_id", zip.ID(), "complex_id", zip.ComplexID())
		return 0, fmt.Errorf("update complex zip code: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.zip_code.update.rows_affected_error", "error", err, "zip_id", zip.ID(), "complex_id", zip.ComplexID())
		return 0, fmt.Errorf("complex zip code rows affected: %w", err)
	}

	return affected, nil
}
