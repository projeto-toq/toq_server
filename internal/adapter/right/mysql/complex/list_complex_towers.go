package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	complexrepoconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	repository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/complex_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) ListComplexTowers(ctx context.Context, tx *sql.Tx, params repository.ListComplexTowersParams) ([]complexmodel.ComplexTowerInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	var builder strings.Builder
	builder.WriteString("SELECT id, complex_id, tower, floors, total_units, units_per_floor FROM complex_towers WHERE 1=1")
	args := make([]any, 0)

	if params.ComplexID > 0 {
		builder.WriteString(" AND complex_id = ?")
		args = append(args, params.ComplexID)
	}

	if params.Tower != "" {
		builder.WriteString(" AND tower LIKE ?")
		args = append(args, fmt.Sprintf("%%%s%%", params.Tower))
	}

	builder.WriteString(" ORDER BY id ASC")

	if params.Limit > 0 {
		builder.WriteString(" LIMIT ?")
		args = append(args, params.Limit)
	}

	if params.Offset > 0 {
		builder.WriteString(" OFFSET ?")
		args = append(args, params.Offset)
	}

	query := builder.String()

	entities, err := ca.Read(ctx, tx, query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.list_towers.read_error", "error", err, "params", params)
		return nil, fmt.Errorf("list complex towers read: %w", err)
	}

	towers := make([]complexmodel.ComplexTowerInterface, 0, len(entities))

	for _, entity := range entities {
		tower, errConv := complexrepoconverters.ComplexTowerEntityToDomain(entity)
		if errConv != nil {
			utils.SetSpanError(ctx, errConv)
			logger.Error("mysql.complex.list_towers.convert_error", "error", errConv)
			return nil, fmt.Errorf("convert complex tower entity: %w", errConv)
		}

		towers = append(towers, tower)
	}

	return towers, nil
}
