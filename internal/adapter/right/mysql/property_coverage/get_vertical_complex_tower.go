package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetVerticalComplexTower returns a tower by id; sql.ErrNoRows when not found.
func (a *PropertyCoverageAdapter) GetVerticalComplexTower(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.VerticalComplexTowerInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	const query = `
        SELECT id, vertical_complex_id, tower, floors, total_units, units_per_floor
        FROM vertical_complex_towers
        WHERE id = ?
        LIMIT 1;
    `

	row := a.QueryRowContext(ctx, tx, "select", query, id)
	var entity propertycoverageentities.VerticalComplexTowerEntity
	if scanErr := row.Scan(&entity.ID, &entity.VerticalComplexID, &entity.Tower, &entity.Floors, &entity.TotalUnits, &entity.UnitsPerFloor); scanErr != nil {
		if scanErr != sql.ErrNoRows {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.property_coverage.tower.get.scan_error", "err", scanErr)
		}
		return nil, scanErr
	}

	return propertycoverageconverters.VerticalComplexTowerEntityToDomain(entity), nil
}
