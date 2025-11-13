package listingservices

// EndUpdateListingInput carries the listing identity reference required to
// finalize the update workflow for a listing. The UUID must point to an active
// listing identity owned by the current user.
type EndUpdateListingInput struct {
	ListingUUID string
}
