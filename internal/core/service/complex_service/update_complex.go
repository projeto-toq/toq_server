package complexservices

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateComplex atualiza os dados de um empreendimento existente.
func (cs *complexService) UpdateComplex(ctx context.Context, input UpdateComplexInput) (complexmodel.ComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err := ensurePositiveID("id", input.ID); err != nil {
		return nil, err
	}

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

	tx, txErr := cs.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("complex.update.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	success := false
	defer func() {
		if !success {
			if rbErr := cs.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("complex.update.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	current, err := cs.complexRepository.GetComplexByID(ctx, tx, input.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("complex.update.get_error", "err", err, "id", input.ID)
		return nil, utils.InternalError("")
	}

	current.SetName(sanitizeString(input.Name))
	current.SetZipCode(normalizedZip)
	current.SetStreet(normalizeOptional(input.Street))
	current.SetNumber(sanitizeString(input.Number))
	current.SetNeighborhood(normalizeOptional(input.Neighborhood))
	current.SetCity(sanitizeString(input.City))
	current.SetState(sanitizeString(input.State))
	current.SetPhoneNumber(normalizeOptional(input.PhoneNumber))
	current.SetSector(input.Sector)
	current.SetMainRegistration(normalizeOptional(input.MainRegistration))
	current.SetPropertyType(input.PropertyType)

	rows, err := cs.complexRepository.UpdateComplex(ctx, tx, current)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("complex.update.repo_error", "err", err, "id", input.ID)
		return nil, utils.InternalError("")
	}

	if rows == 0 {
		return nil, utils.NotFoundError("complex")
	}

	if cmErr := cs.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("complex.update.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	success = true
	return current, nil
}
