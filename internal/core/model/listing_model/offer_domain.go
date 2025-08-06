package listingmodel

type offer struct {
	id int64
}

func (o *offer) ID() int64 {
	return o.id
}

func (o *offer) SetID(id int64) {
	o.id = id
}
