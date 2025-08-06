package listingentity

type EntityGuarantee struct {
	ID        int64
	ListingID int64
	Priority  uint8
	Guarantee uint8
}
