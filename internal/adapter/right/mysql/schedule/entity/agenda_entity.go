package entity

// AgendaEntity mirrors the listing_agendas table.
type AgendaEntity struct {
	ID                uint64
	ListingIdentityID int64
	OwnerID           int64
	Timezone          string
}
