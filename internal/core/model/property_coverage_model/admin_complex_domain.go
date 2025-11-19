package propertycoveragemodel

import (
	"strings"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

type managedComplex struct {
	id               int64
	kind             CoverageKind
	name             string
	zipCode          string
	street           string
	number           string
	neighborhood     string
	city             string
	state            string
	receptionPhone   string
	sector           Sector
	mainRegistration string
	propertyTypes    globalmodel.PropertyType
	sizes            []VerticalComplexSizeInterface
	towers           []VerticalComplexTowerInterface
	zipCodes         []HorizontalComplexZipCodeInterface
}

func (m *managedComplex) ID() int64 { return m.id }

func (m *managedComplex) SetID(id int64) { m.id = id }

func (m *managedComplex) Kind() CoverageKind { return m.kind }

func (m *managedComplex) SetKind(kind CoverageKind) { m.kind = kind }

func (m *managedComplex) Name() string { return m.name }

func (m *managedComplex) SetName(name string) { m.name = strings.TrimSpace(name) }

func (m *managedComplex) ZipCode() string { return m.zipCode }

func (m *managedComplex) SetZipCode(zip string) { m.zipCode = strings.TrimSpace(zip) }

func (m *managedComplex) Street() string { return m.street }

func (m *managedComplex) SetStreet(street string) { m.street = strings.TrimSpace(street) }

func (m *managedComplex) Number() string { return m.number }

func (m *managedComplex) SetNumber(number string) { m.number = strings.TrimSpace(number) }

func (m *managedComplex) Neighborhood() string { return m.neighborhood }

func (m *managedComplex) SetNeighborhood(value string) { m.neighborhood = strings.TrimSpace(value) }

func (m *managedComplex) City() string { return m.city }

func (m *managedComplex) SetCity(city string) { m.city = strings.TrimSpace(city) }

func (m *managedComplex) State() string { return m.state }

func (m *managedComplex) SetState(state string) { m.state = strings.TrimSpace(state) }

func (m *managedComplex) ReceptionPhone() string { return m.receptionPhone }

func (m *managedComplex) SetReceptionPhone(value string) { m.receptionPhone = strings.TrimSpace(value) }

func (m *managedComplex) Sector() Sector { return m.sector }

func (m *managedComplex) SetSector(sector Sector) { m.sector = sector }

func (m *managedComplex) MainRegistration() string { return m.mainRegistration }

func (m *managedComplex) SetMainRegistration(value string) {
	m.mainRegistration = strings.TrimSpace(value)
}

func (m *managedComplex) PropertyTypes() globalmodel.PropertyType { return m.propertyTypes }

func (m *managedComplex) SetPropertyTypes(value globalmodel.PropertyType) { m.propertyTypes = value }

func (m *managedComplex) Sizes() []VerticalComplexSizeInterface { return m.sizes }

func (m *managedComplex) SetSizes(sizes []VerticalComplexSizeInterface) { m.sizes = sizes }

func (m *managedComplex) AddSize(size VerticalComplexSizeInterface) { m.sizes = append(m.sizes, size) }

func (m *managedComplex) Towers() []VerticalComplexTowerInterface { return m.towers }

func (m *managedComplex) SetTowers(towers []VerticalComplexTowerInterface) { m.towers = towers }

func (m *managedComplex) AddTower(tower VerticalComplexTowerInterface) {
	m.towers = append(m.towers, tower)
}

func (m *managedComplex) ZipCodes() []HorizontalComplexZipCodeInterface { return m.zipCodes }

func (m *managedComplex) SetZipCodes(zips []HorizontalComplexZipCodeInterface) { m.zipCodes = zips }

func (m *managedComplex) AddZipCode(zip HorizontalComplexZipCodeInterface) {
	m.zipCodes = append(m.zipCodes, zip)
}

type verticalComplexSize struct {
	id                int64
	verticalComplexID int64
	size              float64
	description       string
}

func (s *verticalComplexSize) ID() int64 { return s.id }

func (s *verticalComplexSize) SetID(id int64) { s.id = id }

func (s *verticalComplexSize) VerticalComplexID() int64 { return s.verticalComplexID }

func (s *verticalComplexSize) SetVerticalComplexID(id int64) { s.verticalComplexID = id }

func (s *verticalComplexSize) Size() float64 { return s.size }

func (s *verticalComplexSize) SetSize(value float64) { s.size = value }

func (s *verticalComplexSize) Description() string { return s.description }

func (s *verticalComplexSize) SetDescription(value string) { s.description = strings.TrimSpace(value) }

type verticalComplexTower struct {
	id                int64
	verticalComplexID int64
	tower             string
	floors            *int
	totalUnits        *int
	unitsPerFloor     *int
}

func (t *verticalComplexTower) ID() int64 { return t.id }

func (t *verticalComplexTower) SetID(id int64) { t.id = id }

func (t *verticalComplexTower) VerticalComplexID() int64 { return t.verticalComplexID }

func (t *verticalComplexTower) SetVerticalComplexID(id int64) { t.verticalComplexID = id }

func (t *verticalComplexTower) Tower() string { return t.tower }

func (t *verticalComplexTower) SetTower(value string) { t.tower = strings.TrimSpace(value) }

func (t *verticalComplexTower) Floors() *int { return t.floors }

func (t *verticalComplexTower) SetFloors(value *int) { t.floors = value }

func (t *verticalComplexTower) TotalUnits() *int { return t.totalUnits }

func (t *verticalComplexTower) SetTotalUnits(value *int) { t.totalUnits = value }

func (t *verticalComplexTower) UnitsPerFloor() *int { return t.unitsPerFloor }

func (t *verticalComplexTower) SetUnitsPerFloor(value *int) { t.unitsPerFloor = value }

type horizontalComplexZipCode struct {
	id                  int64
	horizontalComplexID int64
	zipCode             string
}

func (z *horizontalComplexZipCode) ID() int64 { return z.id }

func (z *horizontalComplexZipCode) SetID(id int64) { z.id = id }

func (z *horizontalComplexZipCode) HorizontalComplexID() int64 { return z.horizontalComplexID }

func (z *horizontalComplexZipCode) SetHorizontalComplexID(id int64) { z.horizontalComplexID = id }

func (z *horizontalComplexZipCode) ZipCode() string { return z.zipCode }

func (z *horizontalComplexZipCode) SetZipCode(value string) { z.zipCode = strings.TrimSpace(value) }
