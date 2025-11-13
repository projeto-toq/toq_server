package listingmodel

type feature struct {
	id               int64
	listingVersionID int64
	feature_id       int64
	quantity         uint8
}

func (f *feature) ID() int64 {
	return f.id
}

func (f *feature) SetID(id int64) {
	f.id = id
}

func (f *feature) ListingID() int64 {
	// Legacy alias for HTTP DTO backward compatibility - returns listingVersionID.
	// Satellite entities now belong to listing_versions, not listing_identities.
	// Safe to remove only when API versioning allows breaking changes.
	return f.listingVersionID
}

func (f *feature) SetListingID(listingID int64) {
	// Legacy alias for HTTP DTO backward compatibility.
	// Safe to remove only when API versioning allows breaking changes.
	f.listingVersionID = listingID
}

func (f *feature) ListingVersionID() int64 {
	return f.listingVersionID
}

func (f *feature) SetListingVersionID(listingVersionID int64) {
	f.listingVersionID = listingVersionID
}

func (f *feature) FeatureID() int64 {
	return f.feature_id
}

func (f *feature) SetFeatureID(feature int64) {
	f.feature_id = feature
}

func (f *feature) Quantity() uint8 {
	return f.quantity
}

func (f *feature) SetQuantity(quantity uint8) {
	f.quantity = quantity
}
