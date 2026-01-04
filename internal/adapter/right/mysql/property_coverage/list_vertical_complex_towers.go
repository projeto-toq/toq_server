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

// ListVerticalComplexTowers lista torres de um complexo vertical; nunca retorna sql.ErrNoRows.
func (a *PropertyCoverageAdapter) ListVerticalComplexTowers(ctx context.Context, tx *sql.Tx, params propertycoveragerepository.ListVerticalComplexTowersParams) ([]propertycoveragemodel.VerticalComplexTowerInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	var builder strings.Builder
	args := []any{params.VerticalComplexID}

	builder.WriteString("SELECT id, vertical_complex_id, tower, floors, total_units, units_per_floor FROM vertical_complex_towers WHERE vertical_complex_id = ?")

	if strings.TrimSpace(params.Tower) != "" {
		builder.WriteString(" AND tower LIKE ?")
		args = append(args, likePattern(params.Tower))
	}

	builder.WriteString(" ORDER BY id DESC")

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
		logger.Error("mysql.property_coverage.tower.list.query_error", "err", qErr)
		return nil, fmt.Errorf("list vertical towers: %w", qErr)
	}
	defer rows.Close()

	result := make([]propertycoveragemodel.VerticalComplexTowerInterface, 0)

	for rows.Next() {
		var entity propertycoverageentities.VerticalComplexTowerEntity
		if scanErr := rows.Scan(&entity.ID, &entity.VerticalComplexID, &entity.Tower, &entity.Floors, &entity.TotalUnits, &entity.UnitsPerFloor); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.property_coverage.tower.list.scan_error", "err", scanErr)
			return nil, fmt.Errorf("scan tower row: %w", scanErr)
		}
		result = append(result, propertycoverageconverters.VerticalComplexTowerEntityToDomain(entity))
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.tower.list.rows_error", "err", err)
		return nil, fmt.Errorf("iterate tower rows: %w", err)
	}

	return result, nil
}

func (a *PropertyCoverageAdapter) listTowers(ctx context.Context, tx *sql.Tx, complexID int64) ([]propertycoveragemodel.VerticalComplexTowerInterface, error) {
	params := propertycoveragerepository.ListVerticalComplexTowersParams{
		VerticalComplexID: complexID,
	}
	return a.ListVerticalComplexTowers(ctx, tx, params)
}
