package mysqlcomplexadapter

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
)

func (ca *ComplexAdapter) CreateComplexTower(ctx context.Context, tx *sql.Tx, tower complexmodel.ComplexTowerInterface) (int64, error) {
	query := `INSERT INTO complex_towers (
		complex_id,
		tower,
		floors,
		total_units,
		units_per_floor
	) VALUES (?, ?, ?, ?, ?);`

	return ca.Create(
		ctx,
		tx,
		query,
		tower.ComplexID(),
		tower.Tower(),
		nullableIntValue(tower.Floors()),
		nullableIntValue(tower.TotalUnits()),
		nullableIntValue(tower.UnitsPerFloor()),
	)
}
