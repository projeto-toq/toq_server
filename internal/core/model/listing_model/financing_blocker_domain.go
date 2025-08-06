package listingmodel

type financingBlocker struct {
	id        int64
	listingID int64
	blocker   FinancingBlocker
}

func (f *financingBlocker) ID() int64 {
	return f.id
}

func (f *financingBlocker) SetID(id int64) {
	f.id = id
}

func (f *financingBlocker) ListingID() int64 {
	return f.listingID
}

func (f *financingBlocker) SetListingID(listingID int64) {
	f.listingID = listingID
}

func (f *financingBlocker) Blocker() FinancingBlocker {
	return f.blocker
}

func (f *financingBlocker) SetBlocker(blocker FinancingBlocker) {
	f.blocker = blocker
}
