package propertycoverageservice

import (
	"context"
	"sort"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	propertycoveragerepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/property_coverage_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListPublicComplexes returns all vertical and horizontal complexes without pagination for public listing flows.
// It excludes standalone coverage entries because they do not represent named complexes.
func (s *propertyCoverageService) ListPublicComplexes(ctx context.Context) ([]propertycoveragemodel.ManagedComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.list_public.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.list_public.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	vertical := propertycoveragemodel.CoverageKindVertical
	horizontal := propertycoveragemodel.CoverageKindHorizontal

	verticals, err := s.repository.ListManagedComplexes(ctx, tx, propertycoveragerepository.ListManagedComplexesParams{Kind: &vertical})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.list_public.vertical_error", "err", err)
		return nil, utils.InternalError("")
	}

	horizontals, err := s.repository.ListManagedComplexes(ctx, tx, propertycoveragerepository.ListManagedComplexesParams{Kind: &horizontal})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.list_public.horizontal_error", "err", err)
		return nil, utils.InternalError("")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.list_public.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true

	merged := append(verticals, horizontals...)
	sort.SliceStable(merged, func(i, j int) bool {
		return merged[i].Name() < merged[j].Name()
	})

	return merged, nil
}
