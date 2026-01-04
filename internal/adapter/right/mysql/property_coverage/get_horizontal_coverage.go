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

// GetHorizontalCoverage retorna a cobertura horizontal para o CEP informado.
// Retorna sql.ErrNoRows quando não há complexo associado; demais erros são marcados no span e logados.
func (a *PropertyCoverageAdapter) GetHorizontalCoverage(ctx context.Context, tx *sql.Tx, zipCode string) (propertycoveragemodel.CoverageInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	const query = `
	        SELECT hc.id, hc.name, hc.main_registration, hc.type
	    FROM horizontal_complexes hc
	    INNER JOIN horizontal_complex_zip_codes hcz ON hcz.horizontal_complex_id = hc.id
	    WHERE hcz.zip_code = ?
	    LIMIT 1
	`

	row := a.QueryRowContext(ctx, tx, "select", query, zipCode)
	var entity propertycoverageentities.HorizontalCoverageEntity
	if scanErr := row.Scan(&entity.ID, &entity.Name, &entity.MainRegistration, &entity.PropertyTypesBitmask); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, scanErr
		}

		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.property_coverage.get_horizontal.scan_error", "zip_code", zipCode, "err", scanErr)
		return nil, fmt.Errorf("get horizontal coverage: %w", scanErr)
	}

	coverage, convErr := propertycoverageconverters.HorizontalEntityToDomain(entity)
	if convErr != nil {
		utils.SetSpanError(ctx, convErr)
		logger.Error("mysql.property_coverage.get_horizontal.convert_error", "zip_code", zipCode, "err", convErr)
		return nil, convErr
	}

	return coverage, nil
}
