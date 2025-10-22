package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/entity"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListSlotsByRange fetches slots within a period for a specific photographer.
func (a *PhotoSessionAdapter) ListSlotsByRange(ctx context.Context, tx *sql.Tx, photographerID uint64, rangeStart, rangeEnd time.Time) ([]photosessionmodel.PhotographerSlotInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		SELECT id, photographer_user_id, slot_date, slot_start, slot_end, period, status, reservation_token, reserved_until, booked_at
		FROM photographer_time_slots
		WHERE photographer_user_id = ?
		  AND slot_start < ?
		  AND slot_end > ?
		ORDER BY slot_start ASC
	`

	rows, execErr := tx.QueryContext(ctx, query, photographerID, rangeEnd, rangeStart)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.list_slots_by_range.query_error", "err", execErr)
		return nil, fmt.Errorf("list photographer slots by range: %w", execErr)
	}
	defer rows.Close()

	slots := make([]photosessionmodel.PhotographerSlotInterface, 0)

	for rows.Next() {
		var (
			entitySlot       entity.SlotEntity
			reservationToken sql.NullString
			reservedUntil    sql.NullTime
			bookedAt         sql.NullTime
		)

		if err = rows.Scan(
			&entitySlot.ID,
			&entitySlot.PhotographerUserID,
			&entitySlot.SlotDate,
			&entitySlot.SlotStart,
			&entitySlot.SlotEnd,
			&entitySlot.Period,
			&entitySlot.Status,
			&reservationToken,
			&reservedUntil,
			&bookedAt,
		); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("mysql.photo_session.list_slots_by_range.scan_error", "err", err)
			return nil, fmt.Errorf("scan photographer slot: %w", err)
		}

		if reservationToken.Valid {
			token := reservationToken.String
			entitySlot.ReservationToken = &token
		}
		if reservedUntil.Valid {
			value := reservedUntil.Time
			entitySlot.ReservedUntil = &value
		}
		if bookedAt.Valid {
			value := bookedAt.Time
			entitySlot.BookedAt = &value
		}

		slots = append(slots, converters.ToSlotModel(entitySlot))
	}

	if err = rows.Err(); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_slots_by_range.rows_error", "err", err)
		return nil, fmt.Errorf("iterate photographer slots: %w", err)
	}

	return slots, nil
}
