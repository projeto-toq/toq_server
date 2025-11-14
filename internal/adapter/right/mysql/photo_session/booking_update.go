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
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToBookingEntity(booking)

	var reason any
	if entity.Reason.Valid {
		reason = entity.Reason.String
	}

	query := `UPDATE photographer_photo_session_bookings
		SET agenda_entry_id = ?, photographer_user_id = ?, listing_identity_id = ?, starts_at = ?, ends_at = ?, status = ?, reason = ?, updated_at = NOW()
		WHERE id = ?`

	result, execErr := a.ExecContext(
		ctx,
		tx,
		"update",
		query,
		entity.AgendaEntryID,
		entity.PhotographerID,
		entity.ListingIdentityID,
		entity.StartsAt,
		entity.EndsAt,
		entity.Status,
		reason,
		entity.ID,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.photo_session.update_booking.exec_error", "booking_id", entity.ID, "err", execErr)
		return fmt.Errorf("update photographer booking: %w", execErr)
	}

	affected, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.photo_session.update_booking.rows_error", "booking_id", entity.ID, "err", rowsErr)
		return fmt.Errorf("rows affected photographer booking: %w", rowsErr)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
