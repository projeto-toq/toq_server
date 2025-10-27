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

func (a *PhotoSessionAdapter) GetBookingBySlotIDForUpdate(ctx context.Context, tx *sql.Tx, slotID uint64) (photosessionmodel.PhotoSessionBookingInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
        SELECT id, slot_id, listing_id, scheduled_start, scheduled_end, status, notes
        FROM photographer_slot_bookings
        WHERE slot_id = ?
        FOR UPDATE
    `

	row := tx.QueryRowContext(ctx, query, slotID)

	var (
		bookingEntity entity.BookingEntity
		notes         sql.NullString
	)

	if err = row.Scan(
		&bookingEntity.ID,
		&bookingEntity.SlotID,
		&bookingEntity.ListingID,
		&bookingEntity.ScheduledStart,
		&bookingEntity.ScheduledEnd,
		&bookingEntity.Status,
		&notes,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.get_booking_by_slot.scan_error", "err", err, "slot_id", slotID)
		return nil, fmt.Errorf("scan photo session booking by slot: %w", err)
	}

	if notes.Valid {
		value := notes.String
		bookingEntity.Notes = &value
	}

	return converters.ToBookingModel(bookingEntity), nil
}
