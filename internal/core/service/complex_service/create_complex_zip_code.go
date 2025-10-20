package complexservices

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateComplexZipCode cria um novo CEP associado a um empreendimento.
func (cs *complexService) CreateComplexZipCode(ctx context.Context, input CreateComplexZipCodeInput) (complexmodel.ComplexZipCodeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if err := ensurePositiveID("complexId", input.ComplexID); err != nil {
		return nil, err
	}

	normalizedZip, err := normalizeAndValidateZip(input.ZipCode)
	if err != nil {
		return nil, err
	}

	tx, txErr := cs.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("complex.zip.create.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	success := false
	defer func() {
		if !success {
			if rbErr := cs.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("complex.zip.create.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	if _, err = cs.complexRepository.GetComplexByID(ctx, tx, input.ComplexID); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("complex.zip.create.parent_error", "err", err, "complex_id", input.ComplexID)
		return nil, utils.InternalError("")
	}

	zipCode := complexmodel.NewComplexZipCode()
	zipCode.SetComplexID(input.ComplexID)
	zipCode.SetZipCode(normalizedZip)

	id, err := cs.complexRepository.CreateComplexZipCode(ctx, tx, zipCode)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("complex.zip.create.repo_error", "err", err, "complex_id", input.ComplexID)
		return nil, utils.InternalError("")
	}

	zipCode.SetID(id)

	if cmErr := cs.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("complex.zip.create.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	success = true
	return zipCode, nil
}
