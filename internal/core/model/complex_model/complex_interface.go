package complexmodel

import globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"

type ComplexInterface interface {
	ID() int64
	SetID(int64)
	Name() string
	SetName(string)
	ZipCode() string
	SetZipCode(string)
	Street() string
	SetStreet(string)
	Number() string
	SetNumber(string)
	Neighborhood() string
	SetNeighborhood(string)
	City() string
	SetCity(string)
	State() string
	SetState(string)
	PhoneNumber() string
	SetPhoneNumber(string)
	Sector() Sector
	SetSector(Sector)
	GetPropertyType() globalmodel.PropertyType
	SetPropertyType(globalmodel.PropertyType)
	MainRegistration() string
	SetMainRegistration(string)
	ComplexSizes() []ComplexSizeInterface
	SetComplexSizes([]ComplexSizeInterface)
	AddComplexSize(ComplexSizeInterface)
	ComplexTowers() []ComplexTowerInterface
	SetComplexTowers([]ComplexTowerInterface)
	AddComplexTower(ComplexTowerInterface)
	ComplexZipCodes() []ComplexZipCodeInterface
	SetComplexZipCodes([]ComplexZipCodeInterface)
	AddComplexZipCode(ComplexZipCodeInterface)
}

func NewComplex() ComplexInterface {
	return &complex{}
}
