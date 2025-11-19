package propertycoverageservice

import (
	"context"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	propertycoveragerepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/property_coverage_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListComplexes returns managed coverage entries according to the provided filters.
func (s *propertyCoverageService) ListComplexes(ctx context.Context, input ListComplexesInput) ([]propertycoveragemodel.ManagedComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	page, limit := sanitizePagination(input.Page, input.Limit)
	offset := (page - 1) * limit

	params := propertycoveragerepository.ListManagedComplexesParams{
		Name:    sanitizeString(input.Name),
		ZipCode: sanitizeString(input.ZipCode),
		Number:  sanitizeString(input.Number),
		City:    sanitizeString(input.City),
		State:   sanitizeString(input.State),
		Limit:   limit,
		Offset:  offset,
		Sector:  input.Sector,
	}

	if input.PropertyType != nil {
		params.PropertyType = input.PropertyType
	}

	if input.Kind != nil {
		params.Kind = input.Kind
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.list.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.list.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	complexes, err := s.repository.ListManagedComplexes(ctx, tx, params)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.list.repo_error", "err", err)
		return nil, utils.InternalError("")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.list.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true
	return complexes, nil
}
