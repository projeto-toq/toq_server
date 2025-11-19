package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"
	"fmt"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetNoComplexCoverage resolves the standalone coverage entry for the zip code.
func (a *PropertyCoverageAdapter) GetNoComplexCoverage(ctx context.Context, tx *sql.Tx, zipCode string) (propertycoveragemodel.CoverageInterface, error) {
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	const query = `
        SELECT ncz.zip_code, ncz.type
        FROM no_complex_zip_codes ncz
        WHERE ncz.zip_code = ?
        LIMIT 1
    `

	row := a.QueryRowContext(ctx, tx, "select", query, zipCode)
	var entity propertycoverageentities.NoComplexCoverageEntity
	if err := row.Scan(&entity.ZipCode, &entity.PropertyTypesBitmask); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.get_no_complex.scan_error", "zip_code", zipCode, "err", err)
		return nil, fmt.Errorf("get no-complex coverage: %w", err)
	}

	coverage, err := propertycoverageconverters.NoComplexEntityToDomain(entity)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.get_no_complex.convert_error", "zip_code", zipCode, "err", err)
		return nil, err
	}

	return coverage, nil
}
