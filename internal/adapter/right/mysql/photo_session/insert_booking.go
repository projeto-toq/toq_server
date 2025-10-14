package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) InsertBooking(ctx context.Context, tx *sql.Tx, booking photosessionmodel.PhotoSessionBookingInterface) (uint64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		INSERT INTO photographer_slot_bookings (slot_id, listing_id, scheduled_start, scheduled_end, status, created_by, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	res, err := tx.ExecContext(ctx, query,
		booking.SlotID(),
		booking.ListingID(),
		booking.ScheduledStart(),
		booking.ScheduledEnd(),
		string(booking.Status()),
		booking.CreatedBy(),
		booking.Notes(),
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.insert_booking.exec_error", "err", err)
		return 0, fmt.Errorf("insert photo session booking: %w", err)
	}

	insertID, err := res.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.insert_booking.last_insert_id_error", "err", err)
		return 0, fmt.Errorf("photo session booking last insert id: %w", err)
	}

	return uint64(insertID), nil
}
