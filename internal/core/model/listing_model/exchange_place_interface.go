package listingmodel

type ExchangePlaceInterface interface {
	ID() int64
	SetID(id int64)
	ListingID() int64
	SetListingID(listingID int64)
	Neighborhood() string
	SetNeighborhood(neighborhood string)
	City() string
	SetCity(city string)
	State() string
	SetState(state string)
}

func NewExchangePlace() ExchangePlaceInterface {
	return &ExchangePlace{}
}
