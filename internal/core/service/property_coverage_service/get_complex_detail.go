package propertycoverageservice

import (
	"context"
	"database/sql"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetComplexDetail retrieves a managed coverage entry by id/kind.
func (s *propertyCoverageService) GetComplexDetail(ctx context.Context, input GetComplexDetailInput) (propertycoveragemodel.ManagedComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	if err := ensurePositiveID("id", input.ID); err != nil {
		return nil, err
	}

	if err := validateCoverageKind(input.Kind); err != nil {
		return nil, err
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.detail.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.detail.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	entity, err := s.repository.GetManagedComplex(ctx, tx, input.ID, input.Kind)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.detail.repo_error", "err", err, "id", input.ID, "kind", input.Kind)
		return nil, utils.InternalError("")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.detail.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true
	return entity, nil
}
