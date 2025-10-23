package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const bulkSlotChunkSize = 100

// BulkUpsertSlots creates or updates photographer slots ensuring uniqueness per photographer/start.
func (a *PhotoSessionAdapter) BulkUpsertSlots(ctx context.Context, tx *sql.Tx, slots []photosessionmodel.PhotographerSlotInterface) error {
	if len(slots) == 0 {
		return nil
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	for start := 0; start < len(slots); start += bulkSlotChunkSize {
		end := start + bulkSlotChunkSize
		if end > len(slots) {
			end = len(slots)
		}

		chunk := slots[start:end]
		placeholders := make([]string, 0, len(chunk))
		args := make([]any, 0, len(chunk)*9)

		for _, slot := range chunk {
			placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")

			var reservationToken any
			if token := slot.ReservationToken(); token != nil {
				reservationToken = *token
			}

			var reservedUntil any
			if until := slot.ReservedUntil(); until != nil {
				reservedUntil = *until
			}

			var bookedAt any
			if booked := slot.BookedAt(); booked != nil {
				bookedAt = *booked
			}

			args = append(args,
				slot.PhotographerUserID(),
				slot.SlotDate(),
				slot.SlotStart(),
				slot.SlotEnd(),
				string(slot.Period()),
				string(slot.Status()),
				reservationToken,
				reservedUntil,
				bookedAt,
			)
		}

		query := fmt.Sprintf(`
			INSERT INTO photographer_time_slots
				(photographer_user_id, slot_date, slot_start, slot_end, period, status, reservation_token, reserved_until, booked_at)
			VALUES %s
			ON DUPLICATE KEY UPDATE
				slot_date = VALUES(slot_date),
				slot_start = VALUES(slot_start),
				slot_end = VALUES(slot_end),
				period = VALUES(period),
				status = VALUES(status),
				reservation_token = VALUES(reservation_token),
				reserved_until = VALUES(reserved_until),
				booked_at = VALUES(booked_at)
		`, strings.Join(placeholders, ","))

		if _, execErr := tx.ExecContext(ctx, query, args...); execErr != nil {
			utils.SetSpanError(ctx, execErr)
			logger.Error("mysql.photo_session.bulk_upsert_slots.exec_error", "err", execErr)
			return fmt.Errorf("bulk upsert photographer slots: %w", execErr)
		}
	}

	return nil
}
