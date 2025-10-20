package complexservices

import (
	"context"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	repository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/complex_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListComplexes retorna empreendimentos com filtros e paginação.
func (cs *complexService) ListComplexes(ctx context.Context, filter ListComplexesInput) ([]complexmodel.ComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	page, limit := sanitizePagination(filter.Page, filter.Limit)
	offset := (page - 1) * limit

	params := repository.ListComplexesParams{
		Name:         sanitizeString(filter.Name),
		ZipCode:      sanitizeString(filter.ZipCode),
		City:         sanitizeString(filter.City),
		State:        sanitizeString(filter.State),
		Sector:       filter.Sector,
		PropertyType: filter.PropertyType,
		Limit:        limit,
		Offset:       offset,
	}

	tx, txErr := cs.gsi.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("complex.list.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}
	success := false
	defer func() {
		if !success {
			if rbErr := cs.gsi.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("complex.list.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	complexes, err := cs.complexRepository.ListComplexes(ctx, tx, params)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("complex.list.repo_error", "err", err)
		return nil, utils.InternalError("")
	}

	if cmErr := cs.gsi.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("complex.list.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	success = true
	return complexes, nil
}
