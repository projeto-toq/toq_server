package mediaprocessingservice

import (
	"github.com/projeto-toq/toq_server/internal/core/derrors"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// ensureProjectFlowAllowed validates property type and status for project media flow.
func (s *mediaProcessingService) ensureProjectFlowAllowed(listing listingmodel.ListingVersionInterface) error {
	if listing == nil {
		return derrors.Validation("listing is required", nil)
	}

	if listing.ListingType() != globalmodel.OffPlanHouse {
		return derrors.Forbidden("project media allowed only for OffPlanHouse listings", nil)
	}

	if listing.Status() != listingmodel.StatusPendingPlanLoading {
		return derrors.Conflict("listing must be in PENDING_PLAN_LOADING to handle project media")
	}

	return nil
}
