package cepmodel

type cep struct {
	cep          string
	street       string
	complement   string
	neighborhood string
	city         string
	state        string
}

func (c *cep) GetCep() string {
	return c.cep
}

func (c *cep) SetCep(cep string) {
	c.cep = cep
}

func (c *cep) GetStreet() string {
	return c.street
}

func (c *cep) SetStreet(street string) {
	c.street = street
}

func (c *cep) GetComplement() string {
	return c.complement
}

func (c *cep) SetComplement(complement string) {
	c.complement = complement
}

func (c *cep) GetNeighborhood() string {
	return c.neighborhood
}

func (c *cep) SetNeighborhood(neighborhood string) {
	c.neighborhood = neighborhood
}

func (c *cep) GetCity() string {
	return c.city
}

func (c *cep) SetCity(city string) {
	c.city = city
}

func (c *cep) GetState() string {
	return c.state
}

func (c *cep) SetState(state string) {
	c.state = state
}
