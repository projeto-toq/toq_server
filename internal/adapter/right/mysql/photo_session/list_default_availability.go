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

// ListDefaultAvailability loads recurring availability slots for a photographer.
func (a *PhotoSessionAdapter) ListDefaultAvailability(ctx context.Context, tx *sql.Tx, photographerID uint64) ([]photosessionmodel.PhotographerDefaultAvailabilityInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		SELECT id, photographer_user_id, weekday, period, start_hour, slots_per_period, slot_duration_minutes
		FROM photographer_default_availability
		WHERE photographer_user_id = ?
		ORDER BY weekday ASC, period ASC
	`

	rows, execErr := tx.QueryContext(ctx, query, photographerID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.list_default_availability.query_error", "err", execErr)
		return nil, fmt.Errorf("list photographer default availability: %w", execErr)
	}
	defer rows.Close()

	records := make([]photosessionmodel.PhotographerDefaultAvailabilityInterface, 0)

	for rows.Next() {
		var e entity.DefaultAvailabilityEntity
		if err = rows.Scan(
			&e.ID,
			&e.PhotographerUserID,
			&e.Weekday,
			&e.Period,
			&e.StartHour,
			&e.SlotsPerPeriod,
			&e.SlotDurationMin,
		); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.photo_session.list_default_availability.scan_error", "err", err)
			return nil, fmt.Errorf("scan photographer default availability: %w", err)
		}

		records = append(records, converters.ToDefaultAvailabilityModel(e))
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_default_availability.rows_error", "err", err)
		return nil, fmt.Errorf("iterate photographer default availability: %w", err)
	}

	return records, nil
}
