package listingservices

// DiscardDraftVersionInput identifies which draft version should be discarded.
type DiscardDraftVersionInput struct {
	// ListingIdentityID references the listing identity (required for ownership validation).
	ListingIdentityID int64
	// VersionID references the draft listing version that must be discarded.
	VersionID int64
}
