package cpfadapter

import (
	"log/slog"
	"time"

	cpfmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cpf_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func ConvertCPFEntityToModel(entity CPFAdapter) (cpf cpfmodel.CPFInterface, err error) {

	cpf = cpfmodel.NewCPF()
	if !entity.Status || entity.Return != "OK" {
		return nil, utils.ErrInternalServer
	}
	cpf.SetNumeroDeCpf(entity.Result.NumeroDeCpf)
	cpf.SetNomeDaPf(entity.Result.NomeDaPf)
	data, err := time.Parse("02/01/2006", entity.Result.DataNascimento)
	if err != nil {
		slog.Error("error converting user born_at to date on validating CPF")
		return nil, utils.ErrInternalServer
	}
	cpf.SetDataNascimento(data)
	cpf.SetSituacaoCadastral(entity.Result.SituacaoCadastral)
	cpf.SetDataInscricao(entity.Result.DataInscricao)
	cpf.SetDigitoVerificador(entity.Result.DigitoVerificador)
	cpf.SetComprovanteEmitido(entity.Result.ComprovanteEmitido)
	cpf.SetComprovanteEmitidoData(entity.Result.ComprovanteEmitidoData)
	return
}
