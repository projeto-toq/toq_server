package listingmodel

type OfferInterface interface {
	ID() int64
	SetID(id int64)
}

func NewOffer() OfferInterface {
	return &offer{}
}
