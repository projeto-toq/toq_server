package cpfadapter

import (
	"fmt"
	"strings"
	"time"

	cpfmodel "github.com/projeto-toq/toq_server/internal/core/model/cpf_model"
	cpfport "github.com/projeto-toq/toq_server/internal/core/port/right/cpf"
)

func ConvertCPFEntityToModel(result cpfResult) (cpfmodel.CPFInterface, error) {
	cpf := cpfmodel.NewCPF()
	cpf.SetNumeroDeCpf(result.NumeroDeCpf)
	cpf.SetNomeDaPf(result.NomeDaPf)

	birthDate, err := time.Parse(cpfDateLayout, result.DataNascimento)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse CPF birth date: %w", cpfport.ErrInfra, err)
	}
	cpf.SetDataNascimento(birthDate)
	cpf.SetSituacaoCadastral(result.SituacaoCadastral)
	cpf.SetDataInscricao(result.DataInscricao)
	cpf.SetDigitoVerificador(result.DigitoVerificador)
	cpf.SetComprovanteEmitido(result.ComprovanteEmitido)
	cpf.SetComprovanteEmitidoData(result.ComprovanteEmitidoData)

	if strings.TrimSpace(result.NumeroDeCpf) == "" {
		return nil, fmt.Errorf("%w: cpf provider returned empty cpf", cpfport.ErrInfra)
	}

	return cpf, nil
}
