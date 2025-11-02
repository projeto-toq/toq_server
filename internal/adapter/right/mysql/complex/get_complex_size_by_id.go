package mysqlcomplexadapter

import (
	"context"
	"database/sql"
	"fmt"

	complexrepoconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/complex/converters"
	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (ca *ComplexAdapter) GetComplexSizeByID(ctx context.Context, tx *sql.Tx, id int64) (complexmodel.ComplexSizeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, complex_id, size, description FROM complex_sizes WHERE id = ? LIMIT 1;`

	rows, err := ca.QueryContext(ctx, tx, "select", query, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_size_by_id.read_error", "error", err, "id", id)
		return nil, fmt.Errorf("get complex size by id query: %w", err)
	}
	defer rows.Close()

	entities, err := rowsToEntities(rows)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_size_by_id.scan_error", "error", err, "id", id)
		return nil, fmt.Errorf("scan complex size rows: %w", err)
	}

	if len(entities) == 0 {
		return nil, sql.ErrNoRows
	}

	size, err := complexrepoconverters.ComplexSizeEntityToDomain(entities[0])
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.complex.get_size_by_id.convert_error", "error", err, "id", id)
		return nil, fmt.Errorf("convert complex size entity: %w", err)
	}

	return size, nil
}
