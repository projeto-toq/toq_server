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

	query := `UPDATE photographer_photo_session_bookings SET status = ? WHERE id = ?`

	result, execErr := a.ExecContext(ctx, tx, "update", query, string(status), bookingID)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.update_booking_status.exec_error", "booking_id", bookingID, "err", execErr)
		return fmt.Errorf("update photographer booking status: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.photo_session.update_booking_status.rows_error", "booking_id", bookingID, "err", rowsErr)
		return fmt.Errorf("rows affected photographer booking status: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
