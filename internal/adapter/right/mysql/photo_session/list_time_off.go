package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListTimeOff retrieves time-off entries intersecting a period for a photographer.
func (a *PhotoSessionAdapter) ListTimeOff(ctx context.Context, tx *sql.Tx, photographerID uint64, rangeStart, rangeEnd time.Time) ([]photosessionmodel.PhotographerTimeOffInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		SELECT id, photographer_user_id, start_date, end_date, reason
		FROM photographer_time_off
		WHERE photographer_user_id = ?
		  AND start_date <= ?
		  AND end_date >= ?
		ORDER BY start_date ASC
	`

	rows, execErr := tx.QueryContext(ctx, query, photographerID, rangeEnd, rangeStart)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.list_time_off.query_error", "err", execErr)
		return nil, fmt.Errorf("list photographer time off: %w", execErr)
	}
	defer rows.Close()

	entries := make([]photosessionmodel.PhotographerTimeOffInterface, 0)

	for rows.Next() {
		var ent entity.TimeOffEntity

		if err = rows.Scan(
			&ent.ID,
			&ent.PhotographerUserID,
			&ent.StartDate,
			&ent.EndDate,
			&ent.Reason,
		); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.photo_session.list_time_off.scan_error", "err", err)
			return nil, fmt.Errorf("scan photographer time off: %w", err)
		}

		entries = append(entries, converters.ToTimeOffModel(ent))
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_time_off.rows_error", "err", err)
		return nil, fmt.Errorf("iterate photographer time off: %w", err)
	}

	return entries, nil
}
