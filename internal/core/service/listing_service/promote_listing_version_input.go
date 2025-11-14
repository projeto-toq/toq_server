package listingservices

// PromoteListingVersionInput carries the version identifier to be promoted to active.
type PromoteListingVersionInput struct {
	// ListingIdentityID references the listing identity (required for ownership validation).
	ListingIdentityID int64
	// VersionID references the listing version that should become the new active version.
	VersionID int64
}
