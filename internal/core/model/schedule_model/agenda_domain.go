package schedulemodel

type agenda struct {
	id        uint64
	listingID int64
	ownerID   int64
	timezone  string
}

func (a *agenda) ID() uint64 {
	return a.id
}

func (a *agenda) SetID(id uint64) {
	a.id = id
}

func (a *agenda) ListingID() int64 {
	return a.listingID
}

func (a *agenda) SetListingID(listingID int64) {
	a.listingID = listingID
}

func (a *agenda) OwnerID() int64 {
	return a.ownerID
}

func (a *agenda) SetOwnerID(ownerID int64) {
	a.ownerID = ownerID
}

func (a *agenda) Timezone() string {
	return a.timezone
}

func (a *agenda) SetTimezone(value string) {
	a.timezone = value
}
