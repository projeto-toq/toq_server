package globalservice

import (
	"context"
	"errors"
	"strings"

	cepmodel "github.com/giulio-alfieri/toq_server/internal/core/model/cep_model"
	cepport "github.com/giulio-alfieri/toq_server/internal/core/port/right/cep"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	validators "github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

func (gs *globalService) GetCEP(ctx context.Context, cep string) (cepmodel.CEPInterface, error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to initialize CEP lookup")
	}
	defer spanEnd()

	addr, err := gs.cep.GetCep(ctx, cep)
	if err != nil {
		switch {
		case errors.Is(err, cepport.ErrInvalid):
			logger.Warn("global_service.get_cep.invalid", "err", err, "cep", maskCEPForLog(cep))
			return nil, utils.ValidationError("zip_code", "Invalid CEP")
		case errors.Is(err, cepport.ErrNotFound):
			logger.Warn("global_service.get_cep.not_found", "err", err, "cep", maskCEPForLog(cep))
			return nil, utils.ValidationError("zip_code", "CEP not found")
		case errors.Is(err, cepport.ErrRateLimited):
			logger.Warn("global_service.get_cep.rate_limited", "err", err, "cep", maskCEPForLog(cep))
			utils.SetSpanError(ctx, err)
			return nil, utils.TooManyAttemptsError("CEP lookup rate limit exceeded")
		case errors.Is(err, cepport.ErrInfra):
			logger.Error("global_service.get_cep.infra", "err", err, "cep", maskCEPForLog(cep))
			utils.SetSpanError(ctx, err)
			return nil, utils.InternalError("Failed to retrieve CEP information")
		}

		logger.Error("global_service.get_cep.unhandled", "err", err, "cep", maskCEPForLog(cep))
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("Failed to retrieve CEP information")
	}

	return addr, nil
}

func maskCEPForLog(raw string) string {
	digits := validators.OnlyDigits(raw)
	length := len(digits)
	if length == 0 {
		return "***"
	}
	visible := 2
	if length < visible {
		visible = 1
	}
	masked := length - visible
	if masked < 3 {
		masked = 3
	}
	return digits[:visible] + strings.Repeat("*", masked)
}
