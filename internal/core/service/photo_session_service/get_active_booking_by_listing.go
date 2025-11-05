package photosessionservices

import (
	"context"
	"database/sql"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

// GetActiveBookingByListingID retrieves the active photo session booking for a listing.
// Returns (nil, sql.ErrNoRows) if no active booking exists.
func (s *photoSessionService) GetActiveBookingByListingID(ctx context.Context, tx *sql.Tx, listingID int64) (photosessionmodel.PhotoSessionBookingInterface, error) {
	// Não inicia tracer pois é método auxiliar chamado dentro de outra transação
	return s.repo.GetActiveBookingByListingID(ctx, tx, listingID)
}
