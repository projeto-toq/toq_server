package listingmodel

type financingBlocker struct {
	id               int64
	listingVersionID int64
	blocker          FinancingBlocker
}

func (f *financingBlocker) ID() int64 {
	return f.id
}

func (f *financingBlocker) SetID(id int64) {
	f.id = id
}

func (f *financingBlocker) ListingID() int64 {
	// Legacy alias for HTTP DTO backward compatibility - returns listingVersionID.
	// Satellite entities now belong to listing_versions, not listing_identities.
	// Safe to remove only when API versioning allows breaking changes.
	return f.listingVersionID
}

func (f *financingBlocker) SetListingID(listingID int64) {
	// Legacy alias for HTTP DTO backward compatibility.
	// Safe to remove only when API versioning allows breaking changes.
	f.listingVersionID = listingID
}

func (f *financingBlocker) ListingVersionID() int64 {
	return f.listingVersionID
}

func (f *financingBlocker) SetListingVersionID(listingVersionID int64) {
	f.listingVersionID = listingVersionID
}

func (f *financingBlocker) Blocker() FinancingBlocker {
	return f.blocker
}

func (f *financingBlocker) SetBlocker(blocker FinancingBlocker) {
	f.blocker = blocker
}
