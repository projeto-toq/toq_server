package mysqlcomplexadapter

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
)

func (ca *ComplexAdapter) UpdateComplexTower(ctx context.Context, tx *sql.Tx, tower complexmodel.ComplexTowerInterface) (int64, error) {
	query := `UPDATE complex_towers SET
		tower = ?,
		floors = ?,
		total_units = ?,
		units_per_floor = ?
	WHERE id = ?
	AND complex_id = ?
	LIMIT 1;`

	return ca.Update(
		ctx,
		tx,
		query,
		tower.Tower(),
		nullableIntValue(tower.Floors()),
		nullableIntValue(tower.TotalUnits()),
		nullableIntValue(tower.UnitsPerFloor()),
		tower.ID(),
		tower.ComplexID(),
	)
}
