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

// GetActiveBookingByListingID retrieves the active booking for a listing if exists.
// Returns sql.ErrNoRows if no active booking is found.
func (a *PhotoSessionAdapter) GetActiveBookingByListingID(ctx context.Context, tx *sql.Tx, listingID int64) (photosessionmodel.PhotoSessionBookingInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Busca bookings com status ativos: PENDING_APPROVAL, ACCEPTED ou ACTIVE
	query := `SELECT id, agenda_entry_id, photographer_user_id, listing_id, starts_at, ends_at, status, reason
		FROM photographer_photo_session_bookings 
		WHERE listing_id = ? 
		AND status IN ('PENDING_APPROVAL', 'ACCEPTED', 'ACTIVE')
		ORDER BY id DESC
		LIMIT 1`

	row := entity.Booking{}
	scanErr := a.QueryRowContext(ctx, tx, "select", query, listingID).Scan(
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
			return nil, sql.ErrNoRows
		}
		utils.SetSpanError(ctx, scanErr)
		logger.Error("mysql.photo_session.get_active_booking_by_listing.scan_error", "listing_id", listingID, "err", scanErr)
		return nil, fmt.Errorf("get active booking by listing: %w", scanErr)
	}

	return converters.ToBookingModel(row), nil
}
