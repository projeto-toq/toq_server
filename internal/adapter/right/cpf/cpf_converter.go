package cpfadapter

import (
	"log/slog"
	"time"

	cpfmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cpf_model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ConvertCPFEntityToModel(entity CPFAdapter) (cpf cpfmodel.CPFInterface, err error) {

	cpf = cpfmodel.NewCPF()
	if !entity.Status || entity.Return != "OK" {
		return nil, status.Error(codes.InvalidArgument, "invalid CPF")
	}
	cpf.SetNumeroDeCpf(entity.Result.NumeroDeCpf)
	cpf.SetNomeDaPf(entity.Result.NomeDaPf)
	data, err := time.Parse("02/01/2006", entity.Result.DataNascimento)
	if err != nil {
		slog.Error("error converting user born_at to date on validating CPF")
		return nil, status.Error(codes.Internal, "internal error")
	}
	cpf.SetDataNascimento(data)
	cpf.SetSituacaoCadastral(entity.Result.SituacaoCadastral)
	cpf.SetDataInscricao(entity.Result.DataInscricao)
	cpf.SetDigitoVerificador(entity.Result.DigitoVerificador)
	cpf.SetComprovanteEmitido(entity.Result.ComprovanteEmitido)
	cpf.SetComprovanteEmitidoData(entity.Result.ComprovanteEmitidoData)
	return
}
