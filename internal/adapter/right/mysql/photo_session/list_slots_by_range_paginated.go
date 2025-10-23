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

// ListSlotsByRangePaginated fetches slots within a period for a specific photographer with pagination support.
func (a *PhotoSessionAdapter) ListSlotsByRangePaginated(ctx context.Context, tx *sql.Tx, photographerID uint64, rangeStart, rangeEnd time.Time, limit, offset int) ([]photosessionmodel.PhotographerSlotInterface, int64, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	countQuery := `
        SELECT COUNT(*)
        FROM photographer_time_slots
        WHERE photographer_user_id = ?
          AND slot_start < ?
          AND slot_end > ?
    `

	var total int64
	if err = tx.QueryRowContext(ctx, countQuery, photographerID, rangeEnd, rangeStart).Scan(&total); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.list_slots_range_paginated.count_error", "err", err)
		return nil, 0, fmt.Errorf("count photographer slots by range: %w", err)
	}

	query := `
        SELECT id, photographer_user_id, slot_date, slot_start, slot_end, period, status, reservation_token, reserved_until, booked_at
        FROM photographer_time_slots
        WHERE photographer_user_id = ?
          AND slot_start < ?
          AND slot_end > ?
        ORDER BY slot_start ASC
        LIMIT ? OFFSET ?
    `

	rows, execErr := tx.QueryContext(ctx, query, photographerID, rangeEnd, rangeStart, limit, offset)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.list_slots_range_paginated.query_error", "err", execErr)
		return nil, 0, fmt.Errorf("list photographer slots by range paginated: %w", execErr)
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
			logger.Error("mysql.photo_session.list_slots_range_paginated.scan_error", "err", err)
			return nil, 0, fmt.Errorf("scan photographer slot: %w", err)
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
		logger.Error("mysql.photo_session.list_slots_range_paginated.rows_error", "err", err)
		return nil, 0, fmt.Errorf("iterate photographer slots: %w", err)
	}

	return slots, total, nil
}
