package listingmodel

type VisitInterface interface {
	ID() int64
	SetID(id int64)
}

func NewVisit() VisitInterface {
	return &visit{}
}
