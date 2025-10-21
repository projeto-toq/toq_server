package holidaymodel

type calendar struct {
	id         uint64
	name       string
	scope      CalendarScope
	state      string
	stateValid bool
	cityIBGE   string
	cityValid  bool
	active     bool
}

func (c *calendar) ID() uint64 {
	return c.id
}

func (c *calendar) SetID(id uint64) {
	c.id = id
}

func (c *calendar) Name() string {
	return c.name
}

func (c *calendar) SetName(value string) {
	c.name = value
}

func (c *calendar) Scope() CalendarScope {
	return c.scope
}

func (c *calendar) SetScope(value CalendarScope) {
	c.scope = value
}

func (c *calendar) State() (string, bool) {
	return c.state, c.stateValid
}

func (c *calendar) SetState(value string) {
	c.state = value
	c.stateValid = true
}

func (c *calendar) ClearState() {
	c.state = ""
	c.stateValid = false
}

func (c *calendar) CityIBGE() (string, bool) {
	return c.cityIBGE, c.cityValid
}

func (c *calendar) SetCityIBGE(value string) {
	c.cityIBGE = value
	c.cityValid = true
}

func (c *calendar) ClearCityIBGE() {
	c.cityIBGE = ""
	c.cityValid = false
}

func (c *calendar) IsActive() bool {
	return c.active
}

func (c *calendar) SetActive(value bool) {
	c.active = value
}
