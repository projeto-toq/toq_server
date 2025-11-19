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

// GetHorizontalCoverage finds the horizontal complex that owns the provided zip code.
func (a *PropertyCoverageAdapter) GetHorizontalCoverage(ctx context.Context, tx *sql.Tx, zipCode string) (propertycoveragemodel.CoverageInterface, error) {
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
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
	if err := row.Scan(&entity.ID, &entity.Name, &entity.MainRegistration, &entity.PropertyTypesBitmask); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.get_horizontal.scan_error", "zip_code", zipCode, "err", err)
		return nil, fmt.Errorf("get horizontal coverage: %w", err)
	}

	coverage, err := propertycoverageconverters.HorizontalEntityToDomain(entity)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.get_horizontal.convert_error", "zip_code", zipCode, "err", err)
		return nil, err
	}

	return coverage, nil
}
