package cnpjmodel

import "time"

type cnpj struct {
	NumeroDeCNPJ   string
	NomeDaPJ       string
	Fantasia       string
	DataNascimento time.Time
}

func (c *cnpj) GetNumeroDeCNPJ() string {
	return c.NumeroDeCNPJ
}

func (c *cnpj) SetNumeroDeCNPJ(numeroCNPJ string) {
	c.NumeroDeCNPJ = numeroCNPJ
}

func (c *cnpj) GetNomeDaPJ() string {
	return c.NomeDaPJ
}

func (c *cnpj) SetNomeDaPJ(nomeDaPJ string) {
	c.NomeDaPJ = nomeDaPJ
}

func (c *cnpj) GetFantasia() string {
	return c.Fantasia
}

func (c *cnpj) SetFantasia(fantasia string) {
	c.Fantasia = fantasia
}

func (c *cnpj) GetDataNascimento() time.Time {
	return c.DataNascimento
}

func (c *cnpj) SetDataNascimento(dataNascimento time.Time) {
	c.DataNascimento = dataNascimento
}
