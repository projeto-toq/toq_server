package complexmodel

type ComplexZipCodeInterface interface {
	ID() int64
	SetID(int64)
	ComplexID() int64
	SetComplexID(int64)
	ZipCode() string
	SetZipCode(string)
}

func NewComplexZipCode() ComplexZipCodeInterface {
	return &complexZipCode{}
}
