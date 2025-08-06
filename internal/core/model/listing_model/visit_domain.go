package listingmodel

type visit struct {
	id int64
}

func (v *visit) ID() int64 {
	return v.id
}

func (v *visit) SetID(id int64) {
	v.id = id
}
