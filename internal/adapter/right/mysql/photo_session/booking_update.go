package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) UpdateBooking(ctx context.Context, tx *sql.Tx, booking photosessionmodel.PhotoSessionBookingInterface) error {
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

	entity := converters.ToBookingEntity(booking)

	var reason any
	if entity.Reason.Valid {
		reason = entity.Reason.String
	}

	query := `UPDATE photographer_photo_session_bookings
		SET agenda_entry_id = ?, photographer_user_id = ?, listing_id = ?, starts_at = ?, ends_at = ?, status = ?, reason = ?
		WHERE id = ?`

	result, err := exec.ExecContext(ctx, query,
		entity.AgendaEntryID,
		entity.PhotographerID,
		entity.ListingID,
		entity.StartsAt,
		entity.EndsAt,
		entity.Status,
		reason,
		entity.ID,
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.update_booking.exec_error", "booking_id", entity.ID, "err", err)
		return fmt.Errorf("update photographer booking: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.update_booking.rows_error", "booking_id", entity.ID, "err", err)
		return fmt.Errorf("rows affected photographer booking: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
