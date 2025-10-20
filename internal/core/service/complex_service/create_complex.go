package complexservices

import (
	"context"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateComplex cria um novo empreendimento com as informações básicas.
func (cs *complexService) CreateComplex(ctx context.Context, input CreateComplexInput) (complexmodel.ComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err := validateRequiredField("name", input.Name); err != nil {
		return nil, err
	}

	normalizedZip, err := normalizeAndValidateZip(input.ZipCode)
	if err != nil {
		return nil, err
	}

	if err := validateRequiredField("number", input.Number); err != nil {
		return nil, err
	}

	if err := validateRequiredField("city", input.City); err != nil {
		return nil, err
	}

	if err := validateRequiredField("state", input.State); err != nil {
		return nil, err
	}

	if err := validateSector(input.Sector); err != nil {
		return nil, err
	}

	domain := complexmodel.NewComplex()
	domain.SetName(sanitizeString(input.Name))
	domain.SetZipCode(normalizedZip)
	domain.SetStreet(normalizeOptional(input.Street))
	domain.SetNumber(sanitizeString(input.Number))
	domain.SetNeighborhood(normalizeOptional(input.Neighborhood))
	domain.SetCity(sanitizeString(input.City))
	domain.SetState(sanitizeString(input.State))
	domain.SetPhoneNumber(normalizeOptional(input.PhoneNumber))
	domain.SetSector(input.Sector)
	domain.SetMainRegistration(normalizeOptional(input.MainRegistration))
	domain.SetPropertyType(input.PropertyType)

	tx, txErr := cs.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("complex.create.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := cs.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("complex.create.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	id, err := cs.complexRepository.CreateComplex(ctx, tx, domain)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("complex.create.repo_error", "err", err)
		return nil, utils.InternalError("")
	}

	domain.SetID(id)

	if cmErr := cs.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("complex.create.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	return domain, nil
}
