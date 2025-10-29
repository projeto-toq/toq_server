package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) ReplaceAssociations(ctx context.Context, tx *sql.Tx, photographerID uint64, calendarIDs []uint64) error {
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

	if _, err := exec.ExecContext(ctx, `DELETE FROM photographer_holiday_calendars WHERE photographer_user_id = ?`, photographerID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.replace_associations.delete_error", "photographer_id", photographerID, "err", err)
		return fmt.Errorf("clear holiday associations: %w", err)
	}

	if len(calendarIDs) == 0 {
		return nil
	}

	stmt, err := exec.PrepareContext(ctx, `INSERT INTO photographer_holiday_calendars (photographer_user_id, holiday_calendar_id) VALUES (?, ?)`)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.replace_associations.prepare_error", "photographer_id", photographerID, "err", err)
		return fmt.Errorf("prepare insert holiday association: %w", err)
	}
	defer stmt.Close()

	for _, calendarID := range calendarIDs {
		if _, err := stmt.ExecContext(ctx, photographerID, calendarID); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.photo_session.replace_associations.insert_error", "photographer_id", photographerID, "calendar_id", calendarID, "err", err)
			return fmt.Errorf("insert holiday association: %w", err)
		}
	}

	return nil
}
