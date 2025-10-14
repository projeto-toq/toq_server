package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) MarkSlotAvailable(ctx context.Context, tx *sql.Tx, slotID uint64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		UPDATE photographer_time_slots
		SET status = 'AVAILABLE', reservation_token = NULL, reserved_until = NULL, booked_at = NULL, updated_at = UTC_TIMESTAMP()
		WHERE id = ?
	`

	res, err := tx.ExecContext(ctx, query, slotID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.slot_available.exec_error", "err", err)
		return fmt.Errorf("set slot available: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.slot_available.rows_affected_error", "err", err)
		return fmt.Errorf("set slot available rows: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
