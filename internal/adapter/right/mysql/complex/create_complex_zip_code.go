package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) CreateComplexZipCode(ctx context.Context, tx *sql.Tx, zip complexmodel.ComplexZipCodeInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO complex_zip_codes (
		complex_id,
		zip_code
	) VALUES (?, ?);`

	result, err := ca.ExecContext(
		ctx,
		tx,
		"insert",
		query,
		zip.ComplexID(),
		zip.ZipCode(),
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.zip_code.create.exec_error", "error", err, "complex_id", zip.ComplexID())
		return 0, fmt.Errorf("insert complex zip code: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.zip_code.create.last_insert_id_error", "error", err, "complex_id", zip.ComplexID())
		return 0, fmt.Errorf("complex zip code last insert id: %w", err)
	}

	return id, nil
}
