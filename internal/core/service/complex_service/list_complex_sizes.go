package complexservices

import (
	"context"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	repository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/complex_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListComplexSizes retorna os tamanhos cadastrados de um empreendimento.
func (cs *complexService) ListComplexSizes(ctx context.Context, filter ListComplexSizesInput) ([]complexmodel.ComplexSizeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if filter.ComplexID > 0 {
		if err := ensurePositiveID("complexId", filter.ComplexID); err != nil {
			return nil, err
		}
	}

	page, limit := sanitizePagination(filter.Page, filter.Limit)
	offset := (page - 1) * limit

	params := repository.ListComplexSizesParams{
		ComplexID: filter.ComplexID,
		Limit:     limit,
		Offset:    offset,
	}

	tx, txErr := cs.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("complex.size.list.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	success := false
	defer func() {
		if !success {
			if rbErr := cs.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("complex.size.list.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	sizes, err := cs.complexRepository.ListComplexSizes(ctx, tx, params)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("complex.size.list.repo_error", "err", err)
		return nil, utils.InternalError("")
	}

	if cmErr := cs.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("complex.size.list.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	success = true
	return sizes, nil
}
