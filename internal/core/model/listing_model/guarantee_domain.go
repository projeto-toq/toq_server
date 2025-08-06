package listingmodel

type guarantee struct {
	id        int64
	listingID int64
	priority  uint8
	guarantee GuaranteeType
}

func (g *guarantee) ID() int64 {
	return g.id
}

func (g *guarantee) SetID(id int64) {
	g.id = id
}

func (g *guarantee) ListingID() int64 {
	return g.listingID
}

func (g *guarantee) SetListingID(listingID int64) {
	g.listingID = listingID
}

func (g *guarantee) Priority() uint8 {
	return g.priority
}

func (g *guarantee) SetPriority(priority uint8) {
	g.priority = priority
}

func (g *guarantee) Guarantee() GuaranteeType {
	return g.guarantee
}

func (g *guarantee) SetGuarantee(guarantee GuaranteeType) {
	g.guarantee = guarantee
}
