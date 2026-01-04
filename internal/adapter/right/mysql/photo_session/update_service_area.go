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

// UpdateServiceArea updates city and state of an existing service area.
func (a *PhotoSessionAdapter) UpdateServiceArea(ctx context.Context, tx *sql.Tx, area photosessionmodel.PhotographerServiceAreaInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	row := converters.ServiceAreaModelToRow(area)
	query := `UPDATE photographer_service_areas SET city = ?, state = ? WHERE id = ?`

	result, execErr := a.ExecContext(ctx, tx, "update", query, strings.TrimSpace(row.City), strings.TrimSpace(row.State), row.ID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.service_area.update.exec_error", "area_id", row.ID, "err", execErr)
		return fmt.Errorf("update service area: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.photo_session.service_area.update.rows_affected_error", "area_id", row.ID, "err", rowsErr)
		return fmt.Errorf("update service area rows affected: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
