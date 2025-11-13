package listingentity

type EntityGuarantee struct {
	ID               int64
	ListingVersionID int64
	Priority         uint8
	Guarantee        uint8
}
