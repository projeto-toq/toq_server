package listingmodel

type GuaranteeInterface interface {
	ID() int64
	SetID(id int64)
	ListingID() int64
	SetListingID(listingID int64)
	ListingVersionID() int64
	SetListingVersionID(listingVersionID int64)
	Priority() uint8
	SetPriority(priority uint8)
	Guarantee() GuaranteeType
	SetGuarantee(guarantee GuaranteeType)
}

func NewGuarantee() GuaranteeInterface {
	return &guarantee{}
}
