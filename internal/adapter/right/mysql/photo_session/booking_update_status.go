package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) UpdateBookingStatus(ctx context.Context, tx *sql.Tx, bookingID uint64, status photosessionmodel.BookingStatus) error {
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

	query := `UPDATE photographer_photo_session_bookings SET status = ? WHERE id = ?`

	result, err := exec.ExecContext(ctx, query, string(status), bookingID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.update_booking_status.exec_error", "booking_id", bookingID, "err", err)
		return fmt.Errorf("update photographer booking status: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.update_booking_status.rows_error", "booking_id", bookingID, "err", err)
		return fmt.Errorf("rows affected photographer booking status: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
