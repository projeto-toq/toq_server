package cnpjmodel

import "time"

type CNPJInterface interface {
	GetNumeroDeCNPJ() string
	SetNumeroDeCNPJ(numeroDeCpf string)
	GetNomeDaPJ() string
	SetNomeDaPJ(nomeDaPf string)
	GetFantasia() string
	SetFantasia(string)
	GetDataNascimento() time.Time
	SetDataNascimento(dataNascimento time.Time)
}

func NewCNPJ() CNPJInterface {
	return &cnpj{}
}
