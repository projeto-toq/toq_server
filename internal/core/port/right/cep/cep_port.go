package cepport

import (
	"context"

	cepmodel "github.com/projeto-toq/toq_server/internal/core/model/cep_model"
)

type CEPPortInterface interface {
	GetCep(ctx context.Context, cepToSearch string) (cep cepmodel.CEPInterface, err error)
}
