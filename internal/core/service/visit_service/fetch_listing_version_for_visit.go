package visitservice

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// fetchListingVersionForVisit loads the listing version referenced by the visit record.
// It ensures the visit detail response has the correct title, description, and address snapshot.
func (s *visitService) fetchListingVersionForVisit(ctx context.Context, tx *sql.Tx, visit listingmodel.VisitInterface) (listingmodel.ListingInterface, error) {
	listing, err := s.listingRepo.GetListingVersionByIdentityAndNumber(ctx, tx, visit.ListingIdentityID(), visit.ListingVersion())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("ListingVersion")
		}

		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("visit.get.fetch_listing_error", "listing_identity_id", visit.ListingIdentityID(), "listing_version", visit.ListingVersion(), "err", err)
		return nil, utils.InternalError("")
	}

	return listing, nil
}
