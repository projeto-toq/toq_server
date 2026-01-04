package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateServiceArea creates a new service area entry.
func (a *PhotoSessionAdapter) CreateServiceArea(ctx context.Context, tx *sql.Tx, area photosessionmodel.PhotographerServiceAreaInterface) (uint64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	row := converters.ServiceAreaModelToRow(area)
	query := `INSERT INTO photographer_service_areas (photographer_user_id, city, state) VALUES (?, ?, ?)`

	execCity := strings.TrimSpace(row.City)
	execState := strings.TrimSpace(row.State)

	result, execErr := a.ExecContext(ctx, tx, "insert", query, row.PhotographerUserID, execCity, execState)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.service_area.create.exec_error", "photographer_id", row.PhotographerUserID, "err", execErr)
		return 0, fmt.Errorf("create service area: %w", execErr)
	}

	insertedID, idErr := result.LastInsertId()
	if idErr != nil {
		utils.SetSpanError(ctx, idErr)
		logger.Error("mysql.photo_session.service_area.create.last_insert_error", "photographer_id", row.PhotographerUserID, "err", idErr)
		return 0, fmt.Errorf("retrieve service area id: %w", idErr)
	}

	return uint64(insertedID), nil
}
