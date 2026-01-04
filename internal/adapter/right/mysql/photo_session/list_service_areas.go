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

// ListAllServiceAreas lists service areas with optional filters.
func (a *PhotoSessionAdapter) ListAllServiceAreas(ctx context.Context, tx *sql.Tx, filter photosessionmodel.ServiceAreaFilter) ([]photosessionmodel.PhotographerServiceAreaInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
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

	rows, queryErr := a.QueryContext(ctx, tx, "select", baseQuery.String(), params...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.photo_session.service_area.list_all.query_error", "err", queryErr)
		return nil, fmt.Errorf("list service areas: %w", queryErr)
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
