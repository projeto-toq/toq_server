package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) MarkSlotBooked(ctx context.Context, tx *sql.Tx, slotID uint64, bookedAt time.Time) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		UPDATE photographer_time_slots
		SET status = 'BOOKED', booked_at = ?, reservation_token = NULL, reserved_until = NULL
		WHERE id = ?
	`

	res, err := tx.ExecContext(ctx, query, bookedAt, slotID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.slot_book.exec_error", "err", err)
		return fmt.Errorf("set slot booked: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.slot_book.rows_affected_error", "err", err)
		return fmt.Errorf("set slot booked rows: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
