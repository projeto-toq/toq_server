package listingservices

// DiscardDraftVersionInput identifies which draft version should be discarded.
type DiscardDraftVersionInput struct {
	// VersionID references the draft listing version that must be discarded.
	VersionID int64
}
