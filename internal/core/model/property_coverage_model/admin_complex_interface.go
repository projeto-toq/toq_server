package propertycoveragemodel

import globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"

// ManagedComplexInterface exposes the metadata required to manage coverage entries
// across vertical, horizontal and standalone sources.
type ManagedComplexInterface interface {
	ID() int64
	SetID(int64)
	Kind() CoverageKind
	SetKind(CoverageKind)
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
	ReceptionPhone() string
	SetReceptionPhone(string)
	Sector() Sector
	SetSector(Sector)
	MainRegistration() string
	SetMainRegistration(string)
	PropertyTypes() globalmodel.PropertyType
	SetPropertyTypes(globalmodel.PropertyType)
	Sizes() []VerticalComplexSizeInterface
	SetSizes([]VerticalComplexSizeInterface)
	AddSize(VerticalComplexSizeInterface)
	Towers() []VerticalComplexTowerInterface
	SetTowers([]VerticalComplexTowerInterface)
	AddTower(VerticalComplexTowerInterface)
	ZipCodes() []HorizontalComplexZipCodeInterface
	SetZipCodes([]HorizontalComplexZipCodeInterface)
	AddZipCode(HorizontalComplexZipCodeInterface)
}

// VerticalComplexSizeInterface describes a single unit size row.
type VerticalComplexSizeInterface interface {
	ID() int64
	SetID(int64)
	VerticalComplexID() int64
	SetVerticalComplexID(int64)
	Size() float64
	SetSize(float64)
	Description() string
	SetDescription(string)
}

// VerticalComplexTowerInterface describes an entry in vertical_complex_towers.
type VerticalComplexTowerInterface interface {
	ID() int64
	SetID(int64)
	VerticalComplexID() int64
	SetVerticalComplexID(int64)
	Tower() string
	SetTower(string)
	Floors() *int
	SetFloors(*int)
	TotalUnits() *int
	SetTotalUnits(*int)
	UnitsPerFloor() *int
	SetUnitsPerFloor(*int)
}

// HorizontalComplexZipCodeInterface represents horizontal CEP coverage rows.
type HorizontalComplexZipCodeInterface interface {
	ID() int64
	SetID(int64)
	HorizontalComplexID() int64
	SetHorizontalComplexID(int64)
	ZipCode() string
	SetZipCode(string)
}

// NewManagedComplex creates an empty managed complex aggregate ready to be filled.
func NewManagedComplex() ManagedComplexInterface {
	return &managedComplex{}
}

// NewVerticalComplexSize creates a new size aggregate.
func NewVerticalComplexSize() VerticalComplexSizeInterface {
	return &verticalComplexSize{}
}

// NewVerticalComplexTower creates a new tower aggregate.
func NewVerticalComplexTower() VerticalComplexTowerInterface {
	return &verticalComplexTower{}
}

// NewHorizontalComplexZipCode creates a new horizontal zip code aggregate.
func NewHorizontalComplexZipCode() HorizontalComplexZipCodeInterface {
	return &horizontalComplexZipCode{}
}
