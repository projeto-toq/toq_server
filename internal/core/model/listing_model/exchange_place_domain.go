package listingmodel

type ExchangePlace struct {
	id               int64
	listingVersionID int64
	neighborhood     string
	city             string
	state            string
}

func (e *ExchangePlace) ID() int64 {
	return e.id
}

func (e *ExchangePlace) SetID(id int64) {
	e.id = id
}

func (e *ExchangePlace) ListingID() int64 {
	// Legacy alias for HTTP DTO backward compatibility - returns listingVersionID.
	// Satellite entities now belong to listing_versions, not listing_identities.
	// Safe to remove only when API versioning allows breaking changes.
	return e.listingVersionID
}

func (e *ExchangePlace) SetListingID(listingID int64) {
	// Legacy alias for HTTP DTO backward compatibility.
	// Safe to remove only when API versioning allows breaking changes.
	e.listingVersionID = listingID
}

func (e *ExchangePlace) ListingVersionID() int64 {
	return e.listingVersionID
}

func (e *ExchangePlace) SetListingVersionID(listingVersionID int64) {
	e.listingVersionID = listingVersionID
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
