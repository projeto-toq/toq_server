package cepadapter

import (
	"fmt"
	"strings"

	cepmodel "github.com/projeto-toq/toq_server/internal/core/model/cep_model"
	cepport "github.com/projeto-toq/toq_server/internal/core/port/right/cep"
)

func ConvertCEPEntityToModel(result cepResult) (cepmodel.CEPInterface, error) {
	normalized := normalizeCEP(result.CEP)
	if len(normalized) != 8 {
		return nil, fmt.Errorf("%w: invalid cep length in response", cepport.ErrInfra)
	}

	cepModel := cepmodel.NewCEP()
	cepModel.SetCep(normalized)
	cepModel.SetStreet(strings.TrimSpace(result.Logradouro))
	cepModel.SetComplement(strings.TrimSpace(result.Complemento))
	cepModel.SetNeighborhood(strings.TrimSpace(result.Bairro))
	cepModel.SetCity(strings.TrimSpace(result.Localidade))
	cepModel.SetState(strings.ToUpper(strings.TrimSpace(result.UF)))

	return cepModel, nil
}
