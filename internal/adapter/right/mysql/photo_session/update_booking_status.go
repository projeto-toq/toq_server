package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) UpdateBookingStatus(ctx context.Context, tx *sql.Tx, bookingID uint64, status photosessionmodel.BookingStatus) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		UPDATE photographer_slot_bookings
		SET status = ?, updated_at = UTC_TIMESTAMP()
		WHERE id = ?
	`

	res, err := tx.ExecContext(ctx, query, string(status), bookingID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.update_booking_status.exec_error", "err", err)
		return fmt.Errorf("update photo session booking status: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.update_booking_status.rows_error", "err", err)
		return fmt.Errorf("update booking status rows: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
