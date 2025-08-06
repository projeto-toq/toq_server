package cpfmodel

import "time"

type CPFInterface interface {
	GetNumeroDeCpf() string
	SetNumeroDeCpf(numeroDeCpf string)
	GetNomeDaPf() string
	SetNomeDaPf(nomeDaPf string)
	GetDataNascimento() time.Time
	SetDataNascimento(dataNascimento time.Time)
	GetSituacaoCadastral() string
	SetSituacaoCadastral(situacaoCadastral string)
	GetDataInscricao() string
	SetDataInscricao(dataInscricao string)
	SetDigitoVerificador(digitoVerificador string)
	GetDigitoVerificador() string
	GetComprovanteEmitido() string
	SetComprovanteEmitido(comprovanteEmitido string)
	GetComprovanteEmitidoData() string
	SetComprovanteEmitidoData(comprovanteEmitidoData string)
}

func NewCPF() CPFInterface {
	return &cpf{}
}
