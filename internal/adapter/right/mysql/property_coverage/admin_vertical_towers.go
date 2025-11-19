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

func (a *PropertyCoverageAdapter) CreateVerticalComplexTower(ctx context.Context, tx *sql.Tx, tower propertycoveragemodel.VerticalComplexTowerInterface) (int64, error) {
	const query = `
		INSERT INTO vertical_complex_towers (
			vertical_complex_id, tower, floors, total_units, units_per_floor
		) VALUES (?, ?, ?, ?, ?);
	`

	args := []any{
		tower.VerticalComplexID(),
		tower.Tower(),
		pointerIntOrZero(tower.Floors()),
		pointerIntOrZero(tower.TotalUnits()),
		pointerIntOrZero(tower.UnitsPerFloor()),
	}

	return a.execInsert(ctx, tx, query, args...)
}

func (a *PropertyCoverageAdapter) UpdateVerticalComplexTower(ctx context.Context, tx *sql.Tx, tower propertycoveragemodel.VerticalComplexTowerInterface) (int64, error) {
	const query = `
		UPDATE vertical_complex_towers SET
			vertical_complex_id = ?,
			tower = ?,
			floors = ?,
			total_units = ?,
			units_per_floor = ?
		WHERE id = ?
		LIMIT 1;
	`

	args := []any{
		tower.VerticalComplexID(),
		tower.Tower(),
		pointerIntOrZero(tower.Floors()),
		pointerIntOrZero(tower.TotalUnits()),
		pointerIntOrZero(tower.UnitsPerFloor()),
		tower.ID(),
	}

	return a.execUpdate(ctx, tx, query, args...)
}

func (a *PropertyCoverageAdapter) DeleteVerticalComplexTower(ctx context.Context, tx *sql.Tx, id int64) (int64, error) {
	const query = "DELETE FROM vertical_complex_towers WHERE id = ? LIMIT 1;"
	return a.execUpdate(ctx, tx, query, id)
}

func (a *PropertyCoverageAdapter) GetVerticalComplexTower(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.VerticalComplexTowerInterface, error) {
	const query = `
		SELECT id, vertical_complex_id, tower, floors, total_units, units_per_floor
		FROM vertical_complex_towers
		WHERE id = ?
		LIMIT 1;
	`

	row := a.QueryRowContext(ctx, tx, "select", query, id)
	var entity propertycoverageentities.VerticalComplexTowerEntity
	if err := row.Scan(&entity.ID, &entity.VerticalComplexID, &entity.Tower, &entity.Floors, &entity.TotalUnits, &entity.UnitsPerFloor); err != nil {
		return nil, err
	}

	return propertycoverageconverters.VerticalComplexTowerEntityToDomain(entity), nil
}

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

	rows, err := a.QueryContext(ctx, tx, "select", builder.String(), args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.property_coverage.tower.list.query_error", "err", err)
		return nil, fmt.Errorf("list vertical towers: %w", err)
	}
	defer rows.Close()

	result := make([]propertycoveragemodel.VerticalComplexTowerInterface, 0)

	for rows.Next() {
		var entity propertycoverageentities.VerticalComplexTowerEntity
		if err := rows.Scan(&entity.ID, &entity.VerticalComplexID, &entity.Tower, &entity.Floors, &entity.TotalUnits, &entity.UnitsPerFloor); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.property_coverage.tower.list.scan_error", "err", err)
			return nil, fmt.Errorf("scan tower row: %w", err)
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
