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

func (a *PropertyCoverageAdapter) CreateVerticalComplexSize(ctx context.Context, tx *sql.Tx, size propertycoveragemodel.VerticalComplexSizeInterface) (int64, error) {
	const query = `
		INSERT INTO vertical_complex_sizes (
			vertical_complex_id, size, description
		) VALUES (?, ?, ?);
	`

	return a.execInsert(ctx, tx, query,
		size.VerticalComplexID(),
		size.Size(),
		nullableString(size.Description()),
	)
}

func (a *PropertyCoverageAdapter) UpdateVerticalComplexSize(ctx context.Context, tx *sql.Tx, size propertycoveragemodel.VerticalComplexSizeInterface) (int64, error) {
	const query = `
		UPDATE vertical_complex_sizes SET
			vertical_complex_id = ?,
			size = ?,
			description = ?
		WHERE id = ?
		LIMIT 1;
	`

	return a.execUpdate(ctx, tx, query,
		size.VerticalComplexID(),
		size.Size(),
		nullableString(size.Description()),
		size.ID(),
	)
}

func (a *PropertyCoverageAdapter) DeleteVerticalComplexSize(ctx context.Context, tx *sql.Tx, id int64) (int64, error) {
	const query = "DELETE FROM vertical_complex_sizes WHERE id = ? LIMIT 1;"
	return a.execUpdate(ctx, tx, query, id)
}

func (a *PropertyCoverageAdapter) GetVerticalComplexSize(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.VerticalComplexSizeInterface, error) {
	const query = `
		SELECT id, vertical_complex_id, size, description
		FROM vertical_complex_sizes
		WHERE id = ?
		LIMIT 1;
	`

	row := a.QueryRowContext(ctx, tx, "select", query, id)
	var entity propertycoverageentities.VerticalComplexSizeEntity
	if err := row.Scan(&entity.ID, &entity.VerticalComplexID, &entity.Size, &entity.Description); err != nil {
		return nil, err
	}

	return propertycoverageconverters.VerticalComplexSizeEntityToDomain(entity), nil
}

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

	rows, err := a.QueryContext(ctx, tx, "select", builder.String(), args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.size.list.query_error", "err", err)
		return nil, fmt.Errorf("list vertical sizes: %w", err)
	}
	defer rows.Close()

	result := make([]propertycoveragemodel.VerticalComplexSizeInterface, 0)

	for rows.Next() {
		var entity propertycoverageentities.VerticalComplexSizeEntity
		if err := rows.Scan(&entity.ID, &entity.VerticalComplexID, &entity.Size, &entity.Description); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.property_coverage.size.list.scan_error", "err", err)
			return nil, fmt.Errorf("scan size row: %w", err)
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
