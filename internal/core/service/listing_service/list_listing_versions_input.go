package listingservices

import listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"

// ListListingVersionsInput filters version history for a listing identity.
type ListListingVersionsInput struct {
	// ListingIdentityID references the identity whose versions should be listed.
	ListingIdentityID int64
	// IncludeDeleted determines whether soft-deleted versions should be returned.
	IncludeDeleted bool
}

// ListListingVersionsOutput aggregates version metadata for presentation layers.
type ListListingVersionsOutput struct {
	Versions []ListingVersionInfo
}

// ListingVersionInfo represents a single version entry along with its active flag.
type ListingVersionInfo struct {
	Version  listingmodel.ListingVersionInterface
	IsActive bool
}
