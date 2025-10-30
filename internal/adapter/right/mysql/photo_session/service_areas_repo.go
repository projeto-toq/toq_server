package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListPhotographerIDsByLocation returns photographer IDs that cover the given city/state.
func (a *PhotoSessionAdapter) ListPhotographerIDsByLocation(ctx context.Context, tx *sql.Tx, city string, state string) ([]uint64, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return nil, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	trimmedCity := strings.TrimSpace(city)
	trimmedState := strings.TrimSpace(state)

	query := `
		SELECT DISTINCT u.id
		FROM users u
		JOIN user_roles ur ON ur.user_id = u.id
		JOIN roles r ON r.id = ur.role_id
		JOIN photographer_service_areas psa ON psa.photographer_user_id = u.id
		WHERE r.slug = 'photographer'
		  AND r.is_active = 1
		  AND ur.is_active = 1
		  AND u.deleted = 0
		  AND psa.city = ?
		  AND psa.state = ?
	`

	rows, err := exec.QueryContext(ctx, query, trimmedCity, trimmedState)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.service_area.list_ids_location.query_error", "city", trimmedCity, "state", trimmedState, "err", err)
		return nil, fmt.Errorf("list photographer ids by location: %w", err)
	}
	defer rows.Close()

	ids := make([]uint64, 0)
	for rows.Next() {
		var id uint64
		if scanErr := rows.Scan(&id); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.photo_session.service_area.list_ids_location.scan_error", "err", scanErr)
			return nil, fmt.Errorf("scan photographer id by location: %w", scanErr)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.service_area.list_ids_location.rows_error", "err", err)
		return nil, fmt.Errorf("iterate photographer ids by location: %w", err)
	}

	return ids, nil
}

// ListServiceAreasByPhotographer lists service areas for a single photographer.
func (a *PhotoSessionAdapter) ListServiceAreasByPhotographer(ctx context.Context, tx *sql.Tx, photographerID uint64) ([]photosessionmodel.PhotographerServiceAreaInterface, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return nil, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, photographer_user_id, city, state FROM photographer_service_areas WHERE photographer_user_id = ? ORDER BY city ASC, state ASC, id ASC`

	rows, err := exec.QueryContext(ctx, query, photographerID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.service_area.list_by_photographer.query_error", "photographer_id", photographerID, "err", err)
		return nil, fmt.Errorf("list service areas by photographer: %w", err)
	}
	defer rows.Close()

	areas := make([]photosessionmodel.PhotographerServiceAreaInterface, 0)
	for rows.Next() {
		var row entity.ServiceArea
		if scanErr := rows.Scan(&row.ID, &row.PhotographerUserID, &row.City, &row.State); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.photo_session.service_area.list_by_photographer.scan_error", "err", scanErr)
			return nil, fmt.Errorf("scan service area by photographer: %w", scanErr)
		}
		areas = append(areas, converters.ServiceAreaRowToModel(row))
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.service_area.list_by_photographer.rows_error", "err", err)
		return nil, fmt.Errorf("iterate service areas by photographer: %w", err)
	}

	return areas, nil
}

// GetServiceAreaByID fetches a single service area.
func (a *PhotoSessionAdapter) GetServiceAreaByID(ctx context.Context, tx *sql.Tx, areaID uint64) (photosessionmodel.PhotographerServiceAreaInterface, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return nil, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, photographer_user_id, city, state FROM photographer_service_areas WHERE id = ?`

	row := exec.QueryRowContext(ctx, query, areaID)

	var entityRow entity.ServiceArea
	if scanErr := row.Scan(&entityRow.ID, &entityRow.PhotographerUserID, &entityRow.City, &entityRow.State); scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.photo_session.service_area.get.query_error", "area_id", areaID, "err", scanErr)
		return nil, fmt.Errorf("get service area by id: %w", scanErr)
	}

	return converters.ServiceAreaRowToModel(entityRow), nil
}

// ListAllServiceAreas lists service areas with optional filters.
func (a *PhotoSessionAdapter) ListAllServiceAreas(ctx context.Context, tx *sql.Tx, filter photosessionmodel.ServiceAreaFilter) ([]photosessionmodel.PhotographerServiceAreaInterface, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return nil, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	baseQuery := strings.Builder{}
	baseQuery.WriteString("SELECT id, photographer_user_id, city, state FROM photographer_service_areas")

	params := make([]any, 0, 4)
	conditions := make([]string, 0, 2)

	if filter.City != nil {
		conditions = append(conditions, "city = ?")
		params = append(params, strings.TrimSpace(*filter.City))
	}
	if filter.State != nil {
		conditions = append(conditions, "state = ?")
		params = append(params, strings.TrimSpace(*filter.State))
	}

	if len(conditions) > 0 {
		baseQuery.WriteString(" WHERE ")
		baseQuery.WriteString(strings.Join(conditions, " AND "))
	}

	baseQuery.WriteString(" ORDER BY city ASC, state ASC, id ASC")

	if filter.Limit > 0 {
		baseQuery.WriteString(" LIMIT ?")
		params = append(params, filter.Limit)
	}
	if filter.Offset > 0 {
		baseQuery.WriteString(" OFFSET ?")
		params = append(params, filter.Offset)
	}

	rows, err := exec.QueryContext(ctx, baseQuery.String(), params...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.service_area.list_all.query_error", "err", err)
		return nil, fmt.Errorf("list service areas: %w", err)
	}
	defer rows.Close()

	areas := make([]photosessionmodel.PhotographerServiceAreaInterface, 0)
	for rows.Next() {
		var row entity.ServiceArea
		if scanErr := rows.Scan(&row.ID, &row.PhotographerUserID, &row.City, &row.State); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.photo_session.service_area.list_all.scan_error", "err", scanErr)
			return nil, fmt.Errorf("scan service area: %w", scanErr)
		}
		areas = append(areas, converters.ServiceAreaRowToModel(row))
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.service_area.list_all.rows_error", "err", err)
		return nil, fmt.Errorf("iterate service areas: %w", err)
	}

	return areas, nil
}

// CreateServiceArea creates a new service area entry.
func (a *PhotoSessionAdapter) CreateServiceArea(ctx context.Context, tx *sql.Tx, area photosessionmodel.PhotographerServiceAreaInterface) (uint64, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return 0, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	row := converters.ServiceAreaModelToRow(area)
	query := `INSERT INTO photographer_service_areas (photographer_user_id, city, state) VALUES (?, ?, ?)`

	execCity := strings.TrimSpace(row.City)
	execState := strings.TrimSpace(row.State)

	result, err := exec.ExecContext(ctx, query, row.PhotographerUserID, execCity, execState)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.service_area.create.exec_error", "photographer_id", row.PhotographerUserID, "err", err)
		return 0, fmt.Errorf("create service area: %w", err)
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.service_area.create.last_insert_error", "photographer_id", row.PhotographerUserID, "err", err)
		return 0, fmt.Errorf("retrieve service area id: %w", err)
	}

	return uint64(insertedID), nil
}

// UpdateServiceArea updates city and state of an existing service area.
func (a *PhotoSessionAdapter) UpdateServiceArea(ctx context.Context, tx *sql.Tx, area photosessionmodel.PhotographerServiceAreaInterface) error {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	row := converters.ServiceAreaModelToRow(area)
	query := `UPDATE photographer_service_areas SET city = ?, state = ? WHERE id = ?`

	result, err := exec.ExecContext(ctx, query, strings.TrimSpace(row.City), strings.TrimSpace(row.State), row.ID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.service_area.update.exec_error", "area_id", row.ID, "err", err)
		return fmt.Errorf("update service area: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.service_area.update.rows_affected_error", "area_id", row.ID, "err", err)
		return fmt.Errorf("update service area rows affected: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteServiceArea removes a service area entry.
func (a *PhotoSessionAdapter) DeleteServiceArea(ctx context.Context, tx *sql.Tx, areaID uint64) error {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `DELETE FROM photographer_service_areas WHERE id = ?`

	result, err := exec.ExecContext(ctx, query, areaID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.service_area.delete.exec_error", "area_id", areaID, "err", err)
		return fmt.Errorf("delete service area: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.service_area.delete.rows_affected_error", "area_id", areaID, "err", err)
		return fmt.Errorf("delete service area rows affected: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
