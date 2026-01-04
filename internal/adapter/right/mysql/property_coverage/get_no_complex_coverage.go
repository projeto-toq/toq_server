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

// GetNoComplexCoverage retorna cobertura standalone para o CEP; sql.ErrNoRows se inexistente.
func (a *PropertyCoverageAdapter) GetNoComplexCoverage(ctx context.Context, tx *sql.Tx, zipCode string) (propertycoveragemodel.CoverageInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
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
	if scanErr := row.Scan(&entity.ZipCode, &entity.PropertyTypesBitmask); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, scanErr
		}

		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.property_coverage.get_no_complex.scan_error", "zip_code", zipCode, "err", scanErr)
		return nil, fmt.Errorf("get no-complex coverage: %w", scanErr)
	}

	coverage, convErr := propertycoverageconverters.NoComplexEntityToDomain(entity)
	if convErr != nil {
		utils.SetSpanError(ctx, convErr)
		logger.Error("mysql.property_coverage.get_no_complex.convert_error", "zip_code", zipCode, "err", convErr)
		return nil, convErr
	}

	return coverage, nil
}
