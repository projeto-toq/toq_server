package listingmodel

type feature struct {
	id         int64
	lisintgID  int64
	feature_id int64
	quantity   uint8
}

func (f *feature) ID() int64 {
	return f.id
}

func (f *feature) SetID(id int64) {
	f.id = id
}

func (f *feature) ListingID() int64 {
	return f.lisintgID
}

func (f *feature) SetListingID(listingID int64) {
	f.lisintgID = listingID
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
