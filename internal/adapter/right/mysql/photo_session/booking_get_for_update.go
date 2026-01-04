package mysqlphotosessionadapter

import (
	"context"
	"database/sql"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// GetBookingByIDForUpdate retrieves a booking with FOR UPDATE locking; requires non-nil transaction for lock to apply.
func (a *PhotoSessionAdapter) GetBookingByIDForUpdate(ctx context.Context, tx *sql.Tx, bookingID uint64) (photosessionmodel.PhotoSessionBookingInterface, error) {
	return a.getBooking(ctx, tx, bookingID, true)
}
