package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) GetBookingForUpdate(ctx context.Context, tx *sql.Tx, bookingID uint64) (photosessionmodel.PhotoSessionBookingInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		SELECT id, slot_id, listing_id, scheduled_start, scheduled_end, status, created_by, notes
		FROM photographer_slot_bookings
		WHERE id = ?
		FOR UPDATE
	`

	row := tx.QueryRowContext(ctx, query, bookingID)

	var (
		entityBooking entity.BookingEntity
		notes         sql.NullString
	)

	if err = row.Scan(
		&entityBooking.ID,
		&entityBooking.SlotID,
		&entityBooking.ListingID,
		&entityBooking.ScheduledStart,
		&entityBooking.ScheduledEnd,
		&entityBooking.Status,
		&entityBooking.CreatedBy,
		&notes,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.get_booking.scan_error", "err", err)
		return nil, fmt.Errorf("scan photo session booking: %w", err)
	}

	if notes.Valid {
		value := notes.String
		entityBooking.Notes = &value
	}

	return converters.ToBookingModel(entityBooking), nil
}
