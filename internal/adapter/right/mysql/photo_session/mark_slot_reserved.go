package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) MarkSlotReserved(ctx context.Context, tx *sql.Tx, slotID uint64, token string, reservedUntil time.Time) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		UPDATE photographer_time_slots
		SET status = 'RESERVED', reservation_token = ?, reserved_until = ?, updated_at = UTC_TIMESTAMP(), booked_at = NULL
		WHERE id = ?
	`

	res, err := tx.ExecContext(ctx, query, token, reservedUntil, slotID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.slot_reserve.exec_error", "err", err)
		return fmt.Errorf("reserve photographer slot: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.slot_reserve.rows_affected_error", "err", err)
		return fmt.Errorf("reserve photographer slot rows: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
