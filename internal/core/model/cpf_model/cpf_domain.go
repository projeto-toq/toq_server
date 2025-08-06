package cpfmodel

import "time"

type cpf struct {
	NumeroDeCpf            string
	NomeDaPf               string
	DataNascimento         time.Time
	SituacaoCadastral      string
	DataInscricao          string
	DigitoVerificador      string
	ComprovanteEmitido     string
	ComprovanteEmitidoData string
}

func (c *cpf) GetNumeroDeCpf() string {
	return c.NumeroDeCpf
}

func (c *cpf) SetNumeroDeCpf(numeroDeCpf string) {
	c.NumeroDeCpf = numeroDeCpf
}

func (c *cpf) GetNomeDaPf() string {
	return c.NomeDaPf
}

func (c *cpf) SetNomeDaPf(nomeDaPf string) {
	c.NomeDaPf = nomeDaPf
}

func (c *cpf) GetDataNascimento() time.Time {
	return c.DataNascimento
}

func (c *cpf) SetDataNascimento(dataNascimento time.Time) {
	c.DataNascimento = dataNascimento
}

func (c *cpf) GetSituacaoCadastral() string {
	return c.SituacaoCadastral
}

func (c *cpf) SetSituacaoCadastral(situacaoCadastral string) {
	c.SituacaoCadastral = situacaoCadastral
}

func (c *cpf) GetDataInscricao() string {
	return c.DataInscricao
}

func (c *cpf) SetDataInscricao(dataInscricao string) {
	c.DataInscricao = dataInscricao
}

func (c *cpf) GetDigitoVerificador() string {
	return c.DigitoVerificador
}

func (c *cpf) SetDigitoVerificador(digitoVerificador string) {
	c.DigitoVerificador = digitoVerificador
}

func (c *cpf) GetComprovanteEmitido() string {
	return c.ComprovanteEmitido
}

func (c *cpf) SetComprovanteEmitido(comprovanteEmitido string) {
	c.ComprovanteEmitido = comprovanteEmitido
}

func (c *cpf) GetComprovanteEmitidoData() string {
	return c.ComprovanteEmitidoData
}

func (c *cpf) SetComprovanteEmitidoData(comprovanteEmitidoData string) {
	c.ComprovanteEmitidoData = comprovanteEmitidoData
}
