package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ReplaceDefaultAvailability wipes existing default availability and inserts the provided set.
func (a *PhotoSessionAdapter) ReplaceDefaultAvailability(ctx context.Context, tx *sql.Tx, photographerID uint64, records []photosessionmodel.PhotographerDefaultAvailabilityInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	deleteQuery := `DELETE FROM photographer_default_availability WHERE photographer_user_id = ?`
	if _, err = tx.ExecContext(ctx, deleteQuery, photographerID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.replace_default_availability.delete_error", "err", err)
		return fmt.Errorf("delete photographer default availability: %w", err)
	}

	if len(records) == 0 {
		return nil
	}

	insertQuery := `
		INSERT INTO photographer_default_availability
			(photographer_user_id, weekday, period, start_hour, slots_per_period, slot_duration_minutes)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	stmt, prepErr := tx.PrepareContext(ctx, insertQuery)
	if prepErr != nil {
		utils.SetSpanError(ctx, prepErr)
		logger.Error("mysql.photo_session.replace_default_availability.prepare_error", "err", prepErr)
		return fmt.Errorf("prepare insert photographer default availability: %w", prepErr)
	}
	defer stmt.Close()

	for _, record := range records {
		entityRecord := converters.FromDefaultAvailabilityModel(record)
		if _, err = stmt.ExecContext(
			ctx,
			entityRecord.PhotographerUserID,
			entityRecord.Weekday,
			entityRecord.Period,
			entityRecord.StartHour,
			entityRecord.SlotsPerPeriod,
			entityRecord.SlotDurationMin,
		); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.photo_session.replace_default_availability.exec_error", "err", err)
			return fmt.Errorf("insert photographer default availability: %w", err)
		}
	}

	return nil
}
