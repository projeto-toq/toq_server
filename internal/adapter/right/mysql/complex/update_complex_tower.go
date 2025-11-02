package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) UpdateComplexTower(ctx context.Context, tx *sql.Tx, tower complexmodel.ComplexTowerInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `UPDATE complex_towers SET
		tower = ?,
		floors = ?,
		total_units = ?,
		units_per_floor = ?
	WHERE id = ?
	AND complex_id = ?
	LIMIT 1;`

	result, err := ca.ExecContext(
		ctx,
		tx,
		"update",
		query,
		tower.Tower(),
		nullableIntValue(tower.Floors()),
		nullableIntValue(tower.TotalUnits()),
		nullableIntValue(tower.UnitsPerFloor()),
		tower.ID(),
		tower.ComplexID(),
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.tower.update.exec_error", "error", err, "tower_id", tower.ID(), "complex_id", tower.ComplexID())
		return 0, fmt.Errorf("update complex tower: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.tower.update.rows_affected_error", "error", err, "tower_id", tower.ID(), "complex_id", tower.ComplexID())
		return 0, fmt.Errorf("complex tower rows affected: %w", err)
	}

	return affected, nil
}
