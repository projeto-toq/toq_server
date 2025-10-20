package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexrepoconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) GetComplexTowerByID(ctx context.Context, tx *sql.Tx, id int64) (complexmodel.ComplexTowerInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, complex_id, tower, floors, total_units, units_per_floor FROM complex_towers WHERE id = ? LIMIT 1;`

	entities, err := ca.Read(ctx, tx, query, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_tower_by_id.read_error", "error", err, "id", id)
		return nil, fmt.Errorf("get complex tower by id read: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	tower, err := complexrepoconverters.ComplexTowerEntityToDomain(entities[0])
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_tower_by_id.convert_error", "error", err, "id", id)
		return nil, fmt.Errorf("convert complex tower entity: %w", err)
	}

	return tower, nil
}
