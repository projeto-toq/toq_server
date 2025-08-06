package listingmodel

type FeatureInterface interface {
	ID() int64
	SetID(id int64)
	ListingID() int64
	SetListingID(listingID int64)
	FeatureID() int64
	SetFeatureID(feature int64)
	Quantity() uint8
	SetQuantity(quantity uint8)
}

func NewFeature() FeatureInterface {
	return &feature{}
}
