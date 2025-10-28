package holidaymodel

type calendar struct {
	id         uint64
	name       string
	scope      CalendarScope
	state      string
	stateValid bool
	city       string
	cityValid  bool
	active     bool
	timezone   string
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

func (c *calendar) City() (string, bool) {
	return c.city, c.cityValid
}

func (c *calendar) SetCity(value string) {
	c.city = value
	c.cityValid = true
}

func (c *calendar) ClearCity() {
	c.city = ""
	c.cityValid = false
}

func (c *calendar) IsActive() bool {
	return c.active
}

func (c *calendar) SetActive(value bool) {
	c.active = value
}

func (c *calendar) Timezone() string {
	return c.timezone
}

func (c *calendar) SetTimezone(value string) {
	c.timezone = value
}
