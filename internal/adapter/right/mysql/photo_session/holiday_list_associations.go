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

func (a *PhotoSessionAdapter) ListAssociations(ctx context.Context, tx *sql.Tx, photographerID uint64) ([]photosessionmodel.HolidayCalendarAssociationInterface, error) {
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

	query := `SELECT id, photographer_user_id, holiday_calendar_id, created_at FROM photographer_holiday_calendars WHERE photographer_user_id = ?`

	rows, err := exec.QueryContext(ctx, query, photographerID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_associations.query_error", "photographer_id", photographerID, "err", err)
		return nil, fmt.Errorf("list photographer holiday associations: %w", err)
	}
	defer rows.Close()

	associations := make([]photosessionmodel.HolidayCalendarAssociationInterface, 0)
	for rows.Next() {
		row := entity.HolidayAssociation{}
		if scanErr := rows.Scan(&row.ID, &row.PhotographerUserID, &row.HolidayCalendarID, &row.CreatedAt); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.photo_session.list_associations.scan_error", "photographer_id", photographerID, "err", scanErr)
			return nil, fmt.Errorf("scan holiday association: %w", scanErr)
		}

		associations = append(associations, converters.ToHolidayAssociationModel(row))
	}

	if err := rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_associations.rows_error", "photographer_id", photographerID, "err", err)
		return nil, fmt.Errorf("iterate holiday associations: %w", err)
	}

	return associations, nil
}
