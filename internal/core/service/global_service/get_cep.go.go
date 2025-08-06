package globalservice

import (
	"context"

	cepmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cep_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (gs *globalService) GetCEP(ctx context.Context, cep string) (address cepmodel.CEPInterface, err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	return gs.cep.GetCep(ctx, cep)

}
