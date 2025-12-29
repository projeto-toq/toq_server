package listingservices

import (
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListingStatusAction describes the allowed owner-triggered transitions for listings.
type ListingStatusAction string

const (
	// ListingStatusActionPublish moves READY listings to PUBLISHED.
	ListingStatusActionPublish ListingStatusAction = "PUBLISH"
	// ListingStatusActionSuspend moves published listings back to READY.
	ListingStatusActionSuspend ListingStatusAction = "SUSPEND"
)

// ChangeListingStatusInput captures all parameters required to execute a status transition.
type ChangeListingStatusInput struct {
	ListingIdentityID int64
	Action            ListingStatusAction
	RequesterUserID   int64
}

// Validate ensures the payload contains consistent data before reaching the repository layer.
func (in ChangeListingStatusInput) Validate() error {
	if in.ListingIdentityID <= 0 {
		return utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}

	switch in.Action {
	case ListingStatusActionPublish, ListingStatusActionSuspend:
	default:
		return utils.ValidationError("action", "action must be either PUBLISH or SUSPEND")
	}

	if in.RequesterUserID <= 0 {
		return utils.AuthorizationError("User context is required")
	}

	return nil
}

// ChangeListingStatusOutput exposes metadata about the applied transition.
type ChangeListingStatusOutput struct {
	ListingIdentityID int64
	ActiveVersionID   int64
	PreviousStatus    listingmodel.ListingStatus
	NewStatus         listingmodel.ListingStatus
}
