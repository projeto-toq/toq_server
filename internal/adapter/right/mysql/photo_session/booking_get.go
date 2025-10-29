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

func (a *PhotoSessionAdapter) GetBookingByID(ctx context.Context, tx *sql.Tx, bookingID uint64) (photosessionmodel.PhotoSessionBookingInterface, error) {
	return a.getBooking(ctx, tx, bookingID, false)
}

func (a *PhotoSessionAdapter) GetBookingByIDForUpdate(ctx context.Context, tx *sql.Tx, bookingID uint64) (photosessionmodel.PhotoSessionBookingInterface, error) {
	return a.getBooking(ctx, tx, bookingID, true)
}

func (a *PhotoSessionAdapter) FindBookingByAgendaEntry(ctx context.Context, tx *sql.Tx, agendaEntryID uint64) (photosessionmodel.PhotoSessionBookingInterface, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return nil, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, agenda_entry_id, photographer_user_id, listing_id, starts_at, ends_at, status, reason
		FROM photographer_photo_session_bookings WHERE agenda_entry_id = ?`

	row := entity.Booking{}
	scanErr := exec.QueryRowContext(ctx, query, agendaEntryID).Scan(
		&row.ID,
		&row.AgendaEntryID,
		&row.PhotographerID,
		&row.ListingID,
		&row.StartsAt,
		&row.EndsAt,
		&row.Status,
		&row.Reason,
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

func (a *PhotoSessionAdapter) getBooking(ctx context.Context, tx *sql.Tx, bookingID uint64, forUpdate bool) (photosessionmodel.PhotoSessionBookingInterface, error) {
	ctx, spanEnd, err := withTracer(ctx)
	if err != nil {
		return nil, err
	}
	if spanEnd != nil {
		defer spanEnd()
	}

	exec := a.executor(tx)
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, agenda_entry_id, photographer_user_id, listing_id, starts_at, ends_at, status, reason
		FROM photographer_photo_session_bookings WHERE id = ?`
	if forUpdate {
		query += " FOR UPDATE"
	}

	row := entity.Booking{}
	scanErr := exec.QueryRowContext(ctx, query, bookingID).Scan(
		&row.ID,
		&row.AgendaEntryID,
		&row.PhotographerID,
		&row.ListingID,
		&row.StartsAt,
		&row.EndsAt,
		&row.Status,
		&row.Reason,
	)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return nil, scanErr
		}
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.photo_session.get_booking.scan_error", "booking_id", bookingID, "err", scanErr)
		return nil, fmt.Errorf("get booking: %w", scanErr)
	}

	return converters.ToBookingModel(row), nil
}
