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

// GetVerticalCoverage busca cobertura vertical por CEP e número; sql.ErrNoRows quando não há match.
func (a *PropertyCoverageAdapter) GetVerticalCoverage(ctx context.Context, tx *sql.Tx, zipCode, number string) (propertycoveragemodel.CoverageInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	const query = `
			SELECT vc.id, vc.name, vc.main_registration, vc.type
			FROM vertical_complexes vc
			WHERE vc.zip_code = ?
			  AND (
					UPPER(REPLACE(TRIM(vc.number), ' ', '')) = ?
					OR FIND_IN_SET(
						?,
						REPLACE(REPLACE(UPPER(TRIM(vc.number)), ' ', ''), ';', ',')
					) > 0
				  )
			LIMIT 1
		`

	row := a.QueryRowContext(ctx, tx, "select", query, zipCode, number, number)
	var entity propertycoverageentities.VerticalCoverageEntity
	if scanErr := row.Scan(&entity.ID, &entity.Name, &entity.MainRegistration, &entity.PropertyTypesBitmask); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, scanErr
		}

		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.property_coverage.get_vertical.scan_error", "zip_code", zipCode, "number", number, "err", scanErr)
		return nil, fmt.Errorf("get vertical coverage: %w", scanErr)
	}

	coverage, convErr := propertycoverageconverters.VerticalEntityToDomain(entity)
	if convErr != nil {
		utils.SetSpanError(ctx, convErr)
		logger.Error("mysql.property_coverage.get_vertical.convert_error", "zip_code", zipCode, "number", number, "err", convErr)
		return nil, convErr
	}

	return coverage, nil
}
