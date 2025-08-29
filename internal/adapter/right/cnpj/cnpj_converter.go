package cnpjadapter

import (
	"time"

	cnpjmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cnpj_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func ConvertCNPJEntityToModel(entity CNPJAdapter) (cnpj cnpjmodel.CNPJInterface, err error) {
	cnpj = cnpjmodel.NewCNPJ()
	if !entity.Status || entity.Return != "OK" {
		return nil, utils.ErrInternalServer
	}
	cnpj.SetNumeroDeCNPJ(entity.Result.NumeroDeCNPJ)
	cnpj.SetNomeDaPJ(entity.Result.NomeDaPJ)
	cnpj.SetFantasia(entity.Result.Fantasia)
	data, erro := time.Parse("02/01/2006", entity.Result.DataNascimento)
	if erro != nil {
		return nil, utils.ErrInternalServer
	}
	cnpj.SetDataNascimento(data)
	return
}
