package listingmodel

type FinancingBlockerInterface interface {
	ID() int64
	SetID(id int64)
	ListingID() int64
	SetListingID(listingID int64)
	Blocker() FinancingBlocker
	SetBlocker(blocker FinancingBlocker)
}

func NewFinancingBlocker() FinancingBlockerInterface {
	return &financingBlocker{}
}
