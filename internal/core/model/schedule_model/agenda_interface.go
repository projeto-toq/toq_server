package schedulemodel

// AgendaInterface exposes the properties of a listing schedule.
type AgendaInterface interface {
	ID() uint64
	SetID(id uint64)
	ListingIdentityID() int64
	SetListingIdentityID(listingIdentityID int64)
	OwnerID() int64
	SetOwnerID(ownerID int64)
	Timezone() string
	SetTimezone(value string)
}

// NewAgenda builds an empty agenda domain object.
func NewAgenda() AgendaInterface {
	return &agenda{}
}
