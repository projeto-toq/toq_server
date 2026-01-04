package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetHorizontalComplexZipCode returns a horizontal complex zip mapping by id; sql.ErrNoRows when missing.
func (a *PropertyCoverageAdapter) GetHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.HorizontalComplexZipCodeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	const query = `
        SELECT id, horizontal_complex_id, zip_code
        FROM horizontal_complex_zip_codes
        WHERE id = ?
        LIMIT 1;
    `

	row := a.QueryRowContext(ctx, tx, "select", query, id)
	var entity propertycoverageentities.HorizontalComplexZipCodeEntity
	if scanErr := row.Scan(&entity.ID, &entity.HorizontalComplexID, &entity.ZipCode); scanErr != nil {
		if scanErr != sql.ErrNoRows {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.property_coverage.zip.get.scan_error", "err", scanErr)
		}
		return nil, scanErr
	}

	return propertycoverageconverters.HorizontalComplexZipCodeEntityToDomain(entity), nil
}
