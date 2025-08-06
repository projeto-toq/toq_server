package cpfport

import (
	"context"
	"time"

	cpfmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cpf_model"
)

type CPFPortInterface interface {
	GetCpf(ctx context.Context, cpfToSearch string, bornAT time.Time) (cpf cpfmodel.CPFInterface, err error)
}
