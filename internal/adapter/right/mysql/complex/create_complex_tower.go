package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) CreateComplexTower(ctx context.Context, tx *sql.Tx, tower complexmodel.ComplexTowerInterface) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `INSERT INTO complex_towers (
		complex_id,
		tower,
		floors,
		total_units,
		units_per_floor
	) VALUES (?, ?, ?, ?, ?);`

	result, err := ca.ExecContext(
		ctx,
		tx,
		"insert",
		query,
		tower.ComplexID(),
		tower.Tower(),
		nullableIntValue(tower.Floors()),
		nullableIntValue(tower.TotalUnits()),
		nullableIntValue(tower.UnitsPerFloor()),
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.tower.create.exec_error", "error", err, "complex_id", tower.ComplexID())
		return 0, fmt.Errorf("insert complex tower: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.tower.create.last_insert_id_error", "error", err, "complex_id", tower.ComplexID())
		return 0, fmt.Errorf("complex tower last insert id: %w", err)
	}

	return id, nil
}
