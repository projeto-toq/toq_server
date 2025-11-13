package listingmodel

type guarantee struct {
	id               int64
	listingVersionID int64
	priority         uint8
	guarantee        GuaranteeType
}

func (g *guarantee) ID() int64 {
	return g.id
}

func (g *guarantee) SetID(id int64) {
	g.id = id
}

func (g *guarantee) ListingID() int64 {
	// Legacy alias for HTTP DTO backward compatibility - returns listingVersionID.
	// Satellite entities now belong to listing_versions, not listing_identities.
	// Safe to remove only when API versioning allows breaking changes.
	return g.listingVersionID
}

func (g *guarantee) SetListingID(listingID int64) {
	// Legacy alias for HTTP DTO backward compatibility.
	// Safe to remove only when API versioning allows breaking changes.
	g.listingVersionID = listingID
}

func (g *guarantee) ListingVersionID() int64 {
	return g.listingVersionID
}

func (g *guarantee) SetListingVersionID(listingVersionID int64) {
	g.listingVersionID = listingVersionID
}

func (g *guarantee) Priority() uint8 {
	return g.priority
}

func (g *guarantee) SetPriority(priority uint8) {
	g.priority = priority
}

func (g *guarantee) Guarantee() GuaranteeType {
	return g.guarantee
}

func (g *guarantee) SetGuarantee(guarantee GuaranteeType) {
	g.guarantee = guarantee
}
