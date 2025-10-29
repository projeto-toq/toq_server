package mysqlphotosessionadapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/photo_session/converters"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (a *PhotoSessionAdapter) CreateBooking(ctx context.Context, tx *sql.Tx, booking photosessionmodel.PhotoSessionBookingInterface) (uint64, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return 0, err
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

	query := `INSERT INTO photographer_photo_session_bookings (
		agenda_entry_id, photographer_user_id, listing_id, starts_at, ends_at, status, reason
	) VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := exec.ExecContext(
		ctx,
		query,
		entity.AgendaEntryID,
		entity.PhotographerID,
		entity.ListingID,
		entity.StartsAt,
		entity.EndsAt,
		entity.Status,
		reason,
	)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.create_booking.exec_error", "agenda_entry_id", entity.AgendaEntryID, "err", err)
		return 0, fmt.Errorf("insert photographer booking: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("mysql.photo_session.create_booking.last_id_error", "agenda_entry_id", entity.AgendaEntryID, "err", err)
		return 0, fmt.Errorf("booking last insert id: %w", err)
	}

	booking.SetID(uint64(id))
	return uint64(id), nil
}
