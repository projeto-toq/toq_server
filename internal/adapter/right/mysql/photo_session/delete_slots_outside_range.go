package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DeleteSlotsOutsideRange removes slots outside the desired scheduling window.
func (a *PhotoSessionAdapter) DeleteSlotsOutsideRange(ctx context.Context, tx *sql.Tx, photographerID uint64, windowStart, windowEnd time.Time) (int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		DELETE FROM photographer_time_slots
		WHERE photographer_user_id = ?
		  AND (slot_end < ? OR slot_start >= ?)
	`

	res, execErr := tx.ExecContext(ctx, query, photographerID, windowStart, windowEnd)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.delete_slots_outside_range.exec_error", "err", execErr)
		return 0, fmt.Errorf("delete photographer slots outside range: %w", execErr)
	}

	rows, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.photo_session.delete_slots_outside_range.rows_error", "err", rowsErr)
		return 0, fmt.Errorf("delete photographer slots rows affected: %w", rowsErr)
	}

	return rows, nil
}
