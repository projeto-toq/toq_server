package cepadapter

import cepmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cep_model"

func ConvertCEPEntityToModel(cep CEPAdapter) (cepModel cepmodel.CEPInterface) {
	cepModel = cepmodel.NewCEP()
	cepModel.SetCep(cep.Result.CEP)
	cepModel.SetStreet(cep.Result.Logradouro)
	cepModel.SetComplement(cep.Result.Complemento)
	cepModel.SetNeighborhood(cep.Result.Bairro)
	cepModel.SetCity(cep.Result.Localidade)
	cepModel.SetState(cep.Result.UF)
	return
}
