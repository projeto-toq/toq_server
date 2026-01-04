package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListServiceAreasByPhotographer lists service areas for a single photographer.
func (a *PhotoSessionAdapter) ListServiceAreasByPhotographer(ctx context.Context, tx *sql.Tx, photographerID uint64) ([]photosessionmodel.PhotographerServiceAreaInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, photographer_user_id, city, state FROM photographer_service_areas WHERE photographer_user_id = ? ORDER BY city ASC, state ASC, id ASC`

	rows, queryErr := a.QueryContext(ctx, tx, "select", query, photographerID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.photo_session.service_area.list_by_photographer.query_error", "photographer_id", photographerID, "err", queryErr)
		return nil, fmt.Errorf("list service areas by photographer: %w", queryErr)
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
