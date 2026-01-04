package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateVerticalComplexTower updates a tower; returns sql.ErrNoRows when no rows are affected.
func (a *PropertyCoverageAdapter) UpdateVerticalComplexTower(ctx context.Context, tx *sql.Tx, tower propertycoveragemodel.VerticalComplexTowerInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	entity := propertycoverageconverters.VerticalComplexTowerDomainToEntity(tower)

	const query = `
        UPDATE vertical_complex_towers SET
            vertical_complex_id = ?,
            tower = ?,
            floors = ?,
            total_units = ?,
            units_per_floor = ?
        WHERE id = ?
        LIMIT 1;
    `

	return a.execUpdate(ctx, tx, "update", query,
		entity.VerticalComplexID,
		entity.Tower,
		entity.Floors,
		entity.TotalUnits,
		entity.UnitsPerFloor,
		entity.ID,
	)
}
