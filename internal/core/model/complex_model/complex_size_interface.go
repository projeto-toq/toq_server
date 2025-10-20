package complexmodel

type ComplexSizeInterface interface {
	ID() int64
	SetID(id int64)
	ComplexID() int64
	SetComplexID(complexID int64)
	Size() float64
	SetSize(size float64)
	Description() string
	SetDescription(description string)
}

func NewComplexSize() ComplexSizeInterface {
	return &complexSize{}
}
