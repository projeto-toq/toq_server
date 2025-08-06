package cnpjadapter

import (
	"time"

	cnpjmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cnpj_model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ConvertCNPJEntityToModel(entity CNPJAdapter) (cnpj cnpjmodel.CNPJInterface, err error) {
	cnpj = cnpjmodel.NewCNPJ()
	if !entity.Status || entity.Return != "OK" {
		return nil, status.Error(codes.InvalidArgument, "CNPJ not found")
	}
	cnpj.SetNumeroDeCNPJ(entity.Result.NumeroDeCNPJ)
	cnpj.SetNomeDaPJ(entity.Result.NomeDaPJ)
	cnpj.SetFantasia(entity.Result.Fantasia)
	data, erro := time.Parse("02/01/2006", entity.Result.DataNascimento)
	if erro != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	cnpj.SetDataNascimento(data)
	return
}
