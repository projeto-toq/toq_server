package photosessionservices

import (
	"context"
	"database/sql"

	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
)

func (s *photoSessionService) GetActiveBookingByListingIdentityID(ctx context.Context, tx *sql.Tx, listingIdentityID int64) (photosessionmodel.PhotoSessionBookingInterface, error) {
	return s.repo.GetActiveBookingByListingIdentityID(ctx, tx, listingIdentityID)
}
