package cnpjadapter

import (
	"fmt"
	"strings"
	"time"

	cnpjmodel "github.com/projeto-toq/toq_server/internal/core/model/cnpj_model"
	cnpjport "github.com/projeto-toq/toq_server/internal/core/port/right/cnpj"
)

func ConvertCNPJEntityToModel(result cnpjResult) (cnpjmodel.CNPJInterface, error) {
	cnpj := cnpjmodel.NewCNPJ()

	if strings.TrimSpace(result.NumeroDeCNPJ) == "" {
		return nil, fmt.Errorf("%w: cnpj provider returned empty cnpj", cnpjport.ErrInfra)
	}

	cnpj.SetNumeroDeCNPJ(result.NumeroDeCNPJ)
	cnpj.SetNomeDaPJ(result.NomeDaPJ)
	cnpj.SetFantasia(result.Fantasia)

	openingDate, err := time.Parse(cnpjDateLayout, result.DataNascimento)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to parse CNPJ opening date: %w", cnpjport.ErrInfra, err)
	}
	cnpj.SetDataNascimento(openingDate)

	return cnpj, nil
}
