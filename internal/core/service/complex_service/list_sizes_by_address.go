package complexservices

import (
	"context"
	"database/sql"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	repository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/complex_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListSizesByAddress retorna os tamanhos disponíveis a partir de um CEP e, opcionalmente, do número.
func (cs *complexService) ListSizesByAddress(ctx context.Context, input ListSizesByAddressInput) ([]complexmodel.ComplexSizeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	normalizedZip, err := normalizeAndValidateZip(input.ZipCode)
	if err != nil {
		return nil, err
	}
	number := sanitizeString(input.Number)

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := cs.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("complex.size_by_address.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	success := false
	defer func() {
		if !success {
			if rbErr := cs.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("complex.size_by_address.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	var complex complexmodel.ComplexInterface

	if number != "" {
		complex, err = cs.complexRepository.GetVerticalComplex(ctx, tx, normalizedZip, number)
		if err != nil {
			if err != sql.ErrNoRows {
				utils.SetSpanError(ctx, err)
				logger.Error("complex.size_by_address.get_vertical_error", "err", err, "zip", normalizedZip, "number", number)
				return nil, utils.InternalError("")
			}
		}
	}

	if complex == nil {
		complex, err = cs.complexRepository.GetHorizontalComplex(ctx, tx, normalizedZip)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, utils.NotFoundError("complex")
			}
			utils.SetSpanError(ctx, err)
			logger.Error("complex.size_by_address.get_horizontal_error", "err", err, "zip", normalizedZip)
			return nil, utils.InternalError("")
		}
	}

	sizes, err := cs.complexRepository.ListComplexSizes(ctx, tx, repository.ListComplexSizesParams{ComplexID: complex.ID()})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("complex.size_by_address.list_sizes_error", "err", err, "complex_id", complex.ID())
		return nil, utils.InternalError("")
	}

	if cmErr := cs.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("complex.size_by_address.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	success = true
	return sizes, nil
}
