package mysqlpropertycoverageadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	propertycoverageconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/converters"
	propertycoverageentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/property_coverage/entities"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	propertycoveragerepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/property_coverage_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListVerticalComplexSizes lista tamanhos de um complexo vertical; nunca retorna sql.ErrNoRows.
func (a *PropertyCoverageAdapter) ListVerticalComplexSizes(ctx context.Context, tx *sql.Tx, params propertycoveragerepository.ListVerticalComplexSizesParams) ([]propertycoveragemodel.VerticalComplexSizeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	var builder strings.Builder
	args := []any{params.VerticalComplexID}

	builder.WriteString("SELECT id, vertical_complex_id, size, description FROM vertical_complex_sizes WHERE vertical_complex_id = ? ORDER BY id DESC")

	if params.Limit > 0 {
		builder.WriteString(" LIMIT ?")
		args = append(args, params.Limit)
	}

	if params.Offset > 0 {
		builder.WriteString(" OFFSET ?")
		args = append(args, params.Offset)
	}

	rows, qErr := a.QueryContext(ctx, tx, "select", builder.String(), args...)
	if qErr != nil {
		utils.SetSpanError(ctx, qErr)
		logger.Error("mysql.property_coverage.size.list.query_error", "err", qErr)
		return nil, fmt.Errorf("list vertical sizes: %w", qErr)
	}
	defer rows.Close()

	result := make([]propertycoveragemodel.VerticalComplexSizeInterface, 0)

	for rows.Next() {
		var entity propertycoverageentities.VerticalComplexSizeEntity
		if scanErr := rows.Scan(&entity.ID, &entity.VerticalComplexID, &entity.Size, &entity.Description); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.property_coverage.size.list.scan_error", "err", scanErr)
			return nil, fmt.Errorf("scan size row: %w", scanErr)
		}
		result = append(result, propertycoverageconverters.VerticalComplexSizeEntityToDomain(entity))
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.size.list.rows_error", "err", err)
		return nil, fmt.Errorf("iterate size rows: %w", err)
	}

	return result, nil
}

func (a *PropertyCoverageAdapter) listSizes(ctx context.Context, tx *sql.Tx, complexID int64) ([]propertycoveragemodel.VerticalComplexSizeInterface, error) {
	params := propertycoveragerepository.ListVerticalComplexSizesParams{
		VerticalComplexID: complexID,
	}
	return a.ListVerticalComplexSizes(ctx, tx, params)
}
