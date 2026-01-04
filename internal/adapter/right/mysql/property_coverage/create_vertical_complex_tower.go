package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateVerticalComplexTower inserts a tower for a vertical complex and returns the created id.
func (a *PropertyCoverageAdapter) CreateVerticalComplexTower(ctx context.Context, tx *sql.Tx, tower propertycoveragemodel.VerticalComplexTowerInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	entity := propertycoverageconverters.VerticalComplexTowerDomainToEntity(tower)

	const query = `
        INSERT INTO vertical_complex_towers (
            vertical_complex_id, tower, floors, total_units, units_per_floor
        ) VALUES (?, ?, ?, ?, ?);
    `

	return a.execInsert(ctx, tx, query,
		entity.VerticalComplexID,
		entity.Tower,
		entity.Floors,
		entity.TotalUnits,
		entity.UnitsPerFloor,
	)
}
