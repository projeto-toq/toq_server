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

// GetServiceAreaByID fetches a single service area.
func (a *PhotoSessionAdapter) GetServiceAreaByID(ctx context.Context, tx *sql.Tx, areaID uint64) (photosessionmodel.PhotographerServiceAreaInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, photographer_user_id, city, state FROM photographer_service_areas WHERE id = ?`

	row := a.QueryRowContext(ctx, tx, "select", query, areaID)

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
