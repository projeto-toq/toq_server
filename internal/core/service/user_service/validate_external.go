package userservices

import (
	"context"
	"fmt"
	"time"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	validators "github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

func (us *userService) ValidateCPF(ctx context.Context, nationalID string, bornAt time.Time) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to initialize validation tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if us.cpf == nil {
		err := fmt.Errorf("cpf adapter not configured")
		logger.Error("user.validate_cpf.adapter_missing", "err", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("CPF validation service unavailable")
	}

	digits := validators.OnlyDigits(nationalID)
	if len(digits) == 0 {
		return utils.ValidationError("nationalID", "Invalid national ID.")
	}

	if bornAt.IsZero() {
		return utils.ValidationError("bornAt", "Invalid birth date.")
	}

	if _, err := us.cpf.GetCpf(ctx, digits, bornAt); err != nil {
		return us.handleCPFValidationError(ctx, "", err)
	}

	return nil
}

func (us *userService) ValidateCNPJ(ctx context.Context, nationalID string) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to initialize validation tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if us.cnpj == nil {
		err := fmt.Errorf("cnpj adapter not configured")
		logger.Error("user.validate_cnpj.adapter_missing", "err", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("CNPJ validation service unavailable")
	}

	digits := validators.OnlyDigits(nationalID)
	if len(digits) == 0 {
		return utils.ValidationError("nationalID", "Invalid national ID.")
	}

	if _, err := us.cnpj.GetCNPJ(ctx, digits); err != nil {
		return us.handleCNPJValidationError(ctx, "", err)
	}

	return nil
}
