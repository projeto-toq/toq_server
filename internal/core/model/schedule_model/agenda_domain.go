package schedulemodel

type agenda struct {
	id                uint64
	listingIdentityID int64
	ownerID           int64
	timezone          string
}

func (a *agenda) ID() uint64 {
	return a.id
}

func (a *agenda) SetID(id uint64) {
	a.id = id
}

func (a *agenda) ListingIdentityID() int64 {
	return a.listingIdentityID
}

func (a *agenda) SetListingIdentityID(listingIdentityID int64) {
	a.listingIdentityID = listingIdentityID
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
