package propertycoverageservice

import (
	"context"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	propertycoveragerepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/property_coverage_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListComplexTowers lists tower entries for a vertical complex.
func (s *propertyCoverageService) ListComplexTowers(ctx context.Context, input ListComplexTowersInput) ([]propertycoveragemodel.VerticalComplexTowerInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	if err := ensurePositiveID("verticalComplexId", input.VerticalComplexID); err != nil {
		return nil, err
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	page, limit := sanitizePagination(input.Page, input.Limit)
	offset := (page - 1) * limit

	params := propertycoveragerepository.ListVerticalComplexTowersParams{
		VerticalComplexID: input.VerticalComplexID,
		Tower:             sanitizeString(input.Tower),
		Limit:             limit,
		Offset:            offset,
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.tower.list.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.tower.list.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	towers, err := s.repository.ListVerticalComplexTowers(ctx, tx, params)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.tower.list.repo_error", "err", err)
		return nil, utils.InternalError("")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.tower.list.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true
	return towers, nil
}
