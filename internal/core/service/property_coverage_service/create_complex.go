package propertycoverageservice

import (
	"context"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateComplex registers a managed coverage entry mapped to the new property coverage schema.
func (s *propertyCoverageService) CreateComplex(ctx context.Context, input CreateComplexInput) (propertycoveragemodel.ManagedComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	domain, err := buildManagedComplexFromCreateInput(input)
	if err != nil {
		return nil, err
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.create.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.create.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	id, err := s.repository.CreateManagedComplex(ctx, tx, domain)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.create.repo_error", "err", err)
		return nil, utils.InternalError("")
	}

	domain.SetID(id)

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.create.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true
	return domain, nil
}
