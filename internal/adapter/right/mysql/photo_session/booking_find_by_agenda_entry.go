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

// FindBookingByAgendaEntry fetches a booking using the agenda entry id; returns sql.ErrNoRows when absent.
func (a *PhotoSessionAdapter) FindBookingByAgendaEntry(ctx context.Context, tx *sql.Tx, agendaEntryID uint64) (photosessionmodel.PhotoSessionBookingInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, agenda_entry_id, photographer_user_id, listing_identity_id, starts_at, ends_at, status, reason, reservation_token, reserved_until
		FROM photographer_photo_session_bookings WHERE agenda_entry_id = ?`

	row := entity.Booking{}
	scanErr := a.QueryRowContext(ctx, tx, "select", query, agendaEntryID).Scan(
		&row.ID,
		&row.AgendaEntryID,
		&row.PhotographerID,
		&row.ListingIdentityID,
		&row.StartsAt,
		&row.EndsAt,
		&row.Status,
		&row.Reason,
		&row.ReservationToken,
		&row.ReservedUntil,
	)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, scanErr
		}
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.photo_session.find_booking.scan_error", "agenda_entry_id", agendaEntryID, "err", scanErr)
		return nil, fmt.Errorf("find booking by agenda entry: %w", scanErr)
	}

	return converters.ToBookingModel(row), nil
}
