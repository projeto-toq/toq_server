package cnpjport

import (
	"context"

	cnpjmodel "github.com/projeto-toq/toq_server/internal/core/model/cnpj_model"
)

type CNPJPortInterface interface {
	GetCNPJ(ctx context.Context, cnpjToSearch string) (cnpj cnpjmodel.CNPJInterface, err error)
}
