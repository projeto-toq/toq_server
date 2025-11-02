package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexrepoconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) GetComplexZipCodeByID(ctx context.Context, tx *sql.Tx, id int64) (complexmodel.ComplexZipCodeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, complex_id, zip_code FROM complex_zip_codes WHERE id = ? LIMIT 1;`

	rows, err := ca.QueryContext(ctx, tx, "select", query, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_zip_code_by_id.read_error", "error", err, "id", id)
		return nil, fmt.Errorf("get complex zip code by id query: %w", err)
	}
	defer rows.Close()

	entities, err := rowsToEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_zip_code_by_id.scan_error", "error", err, "id", id)
		return nil, fmt.Errorf("scan complex zip code rows: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	zipCode, err := complexrepoconverters.ComplexZipCodeEntityToDomain(entities[0])
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_zip_code_by_id.convert_error", "error", err, "id", id)
		return nil, fmt.Errorf("convert complex zip code entity: %w", err)
	}

	return zipCode, nil
}
