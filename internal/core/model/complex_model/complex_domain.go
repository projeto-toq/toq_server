package complexmodel

import globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"

type complex struct {
	id               int64
	name             string
	zipCode          string
	street           string
	number           string
	neighborhood     string
	city             string
	state            string
	phoneNumber      string
	sector           Sector
	mainRegistration string
	PropertyType     globalmodel.PropertyType
	complexSizes     []ComplexSizeInterface
	complexTowers    []ComplexTowerInterface
	complexZipCodes  []ComplexZipCodeInterface
}

func (c *complex) ID() int64 {
	return c.id
}

func (c *complex) SetID(id int64) {
	c.id = id
}

func (c *complex) Name() string {
	return c.name
}

func (c *complex) SetName(name string) {
	c.name = name
}

func (c *complex) ZipCode() string {
	return c.zipCode
}

func (c *complex) SetZipCode(zipCode string) {
	c.zipCode = zipCode
}

func (c *complex) Street() string {
	return c.street
}

func (c *complex) SetStreet(street string) {
	c.street = street
}

func (c *complex) Number() string {
	return c.number
}

func (c *complex) SetNumber(number string) {
	c.number = number
}

func (c *complex) Neighborhood() string {
	return c.neighborhood
}

func (c *complex) SetNeighborhood(neighborhood string) {
	c.neighborhood = neighborhood
}

func (c *complex) City() string {
	return c.city
}

func (c *complex) SetCity(city string) {
	c.city = city
}

func (c *complex) State() string {
	return c.state
}

func (c *complex) SetState(state string) {
	c.state = state
}

func (c *complex) PhoneNumber() string {
	return c.phoneNumber
}

func (c *complex) SetPhoneNumber(phoneNumber string) {
	c.phoneNumber = phoneNumber
}

func (c *complex) Sector() Sector {
	return c.sector
}

func (c *complex) SetSector(sector Sector) {
	c.sector = sector
}

func (c *complex) MainRegistration() string {
	return c.mainRegistration
}

func (c *complex) SetMainRegistration(mainRegistration string) {
	c.mainRegistration = mainRegistration
}

func (c *complex) GetPropertyType() globalmodel.PropertyType {
	return c.PropertyType
}

func (c *complex) SetPropertyType(propertyType globalmodel.PropertyType) {
	c.PropertyType = propertyType
}

func (c *complex) ComplexSizes() []ComplexSizeInterface {
	return c.complexSizes
}

func (c *complex) SetComplexSizes(complexSizes []ComplexSizeInterface) {
	c.complexSizes = complexSizes
}

func (c *complex) AddComplexSize(complexSize ComplexSizeInterface) {
	c.complexSizes = append(c.complexSizes, complexSize)
}

func (c *complex) ComplexTowers() []ComplexTowerInterface {
	return c.complexTowers
}

func (c *complex) SetComplexTowers(complexTowers []ComplexTowerInterface) {
	c.complexTowers = complexTowers
}

func (c *complex) AddComplexTower(complexTower ComplexTowerInterface) {
	c.complexTowers = append(c.complexTowers, complexTower)
}

func (c *complex) ComplexZipCodes() []ComplexZipCodeInterface {
	return c.complexZipCodes
}

func (c *complex) SetComplexZipCodes(complexZipCodes []ComplexZipCodeInterface) {
	c.complexZipCodes = complexZipCodes
}

func (c *complex) AddComplexZipCode(complexZipCode ComplexZipCodeInterface) {
	c.complexZipCodes = append(c.complexZipCodes, complexZipCode)
}
