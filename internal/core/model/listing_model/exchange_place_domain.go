package listingmodel

type ExchangePlace struct {
	id           int64
	listingID    int64
	neighborhood string
	city         string
	state        string
}

func (e *ExchangePlace) ID() int64 {
	return e.id
}

func (e *ExchangePlace) SetID(id int64) {
	e.id = id
}

func (e *ExchangePlace) ListingID() int64 {
	return e.listingID
}

func (e *ExchangePlace) SetListingID(listingID int64) {
	e.listingID = listingID
}

func (e *ExchangePlace) Neighborhood() string {
	return e.neighborhood
}

func (e *ExchangePlace) SetNeighborhood(neighborhood string) {
	e.neighborhood = neighborhood
}

func (e *ExchangePlace) City() string {
	return e.city
}

func (e *ExchangePlace) SetCity(city string) {
	e.city = city
}

func (e *ExchangePlace) State() string {
	return e.state
}

func (e *ExchangePlace) SetState(state string) {
	e.state = state
}
