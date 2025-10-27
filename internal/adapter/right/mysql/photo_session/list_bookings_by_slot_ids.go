package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) ListBookingsBySlotIDs(ctx context.Context, tx *sql.Tx, slotIDs []uint64) ([]photosessionmodel.PhotoSessionBookingInterface, error) {
	if len(slotIDs) == 0 {
		return nil, nil
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	placeholders := make([]string, len(slotIDs))
	args := make([]interface{}, len(slotIDs))
	for i, id := range slotIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
        SELECT id, slot_id, listing_id, scheduled_start, scheduled_end, status, notes
        FROM photographer_slot_bookings
        WHERE slot_id IN (%s)
    `, strings.Join(placeholders, ","))

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_bookings_by_slots.query_error", "err", err)
		return nil, fmt.Errorf("list photo session bookings by slot ids: %w", err)
	}
	defer rows.Close()

	bookings := make([]photosessionmodel.PhotoSessionBookingInterface, 0)

	for rows.Next() {
		var (
			bookingEntity entity.BookingEntity
			notes         sql.NullString
		)

		if err = rows.Scan(
			&bookingEntity.ID,
			&bookingEntity.SlotID,
			&bookingEntity.ListingID,
			&bookingEntity.ScheduledStart,
			&bookingEntity.ScheduledEnd,
			&bookingEntity.Status,
			&notes,
		); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.photo_session.list_bookings_by_slots.scan_error", "err", err)
			return nil, fmt.Errorf("scan photo session booking: %w", err)
		}

		if notes.Valid {
			value := notes.String
			bookingEntity.Notes = &value
		}

		bookings = append(bookings, converters.ToBookingModel(bookingEntity))
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_bookings_by_slots.rows_error", "err", err)
		return nil, fmt.Errorf("iterate photo session bookings: %w", err)
	}

	return bookings, nil
}
