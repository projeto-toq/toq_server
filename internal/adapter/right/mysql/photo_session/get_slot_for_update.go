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

func (a *PhotoSessionAdapter) GetSlotForUpdate(ctx context.Context, tx *sql.Tx, slotID uint64) (photosessionmodel.PhotographerSlotInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `
		SELECT id, photographer_user_id, slot_date, period, status, reservation_token, reserved_until, booked_at
		FROM photographer_time_slots
		WHERE id = ?
		FOR UPDATE
	`

	row := tx.QueryRowContext(ctx, query, slotID)

	var (
		entSlot          entity.SlotEntity
		reservationToken sql.NullString
		reservedUntil    sql.NullTime
		bookedAt         sql.NullTime
	)

	if err = row.Scan(
		&entSlot.ID,
		&entSlot.PhotographerUserID,
		&entSlot.SlotDate,
		&entSlot.Period,
		&entSlot.Status,
		&reservationToken,
		&reservedUntil,
		&bookedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.get_slot.scan_error", "err", err)
		return nil, fmt.Errorf("scan photographer slot: %w", err)
	}

	if reservationToken.Valid {
		token := reservationToken.String
		entSlot.ReservationToken = &token
	}
	if reservedUntil.Valid {
		value := reservedUntil.Time
		entSlot.ReservedUntil = &value
	}
	if bookedAt.Valid {
		value := bookedAt.Time
		entSlot.BookedAt = &value
	}

	return converters.ToSlotModel(entSlot), nil
}
