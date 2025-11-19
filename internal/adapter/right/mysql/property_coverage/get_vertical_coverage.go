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

// GetVerticalCoverage looks for a vertical complex matching the provided zip code and number.
func (a *PropertyCoverageAdapter) GetVerticalCoverage(ctx context.Context, tx *sql.Tx, zipCode, number string) (propertycoveragemodel.CoverageInterface, error) {
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
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
	if err := row.Scan(&entity.ID, &entity.Name, &entity.MainRegistration, &entity.PropertyTypesBitmask); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.get_vertical.scan_error", "zip_code", zipCode, "number", number, "err", err)
		return nil, fmt.Errorf("get vertical coverage: %w", err)
	}

	coverage, err := propertycoverageconverters.VerticalEntityToDomain(entity)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.get_vertical.convert_error", "zip_code", zipCode, "number", number, "err", err)
		return nil, err
	}

	return coverage, nil
}
