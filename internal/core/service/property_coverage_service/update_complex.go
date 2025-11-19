package propertycoverageservice

import (
	"context"
	"database/sql"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateComplex mutates a managed coverage entry based on its coverage kind.
func (s *propertyCoverageService) UpdateComplex(ctx context.Context, input UpdateComplexInput) (propertycoveragemodel.ManagedComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	if err := ensurePositiveID("id", input.ID); err != nil {
		return nil, err
	}

	domain, err := buildManagedComplexFromCreateInput(input.CreateComplexInput)
	if err != nil {
		return nil, err
	}
	domain.SetID(input.ID)

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.update.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.update.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	// Ensure the entity exists and the coverage kind matches before updating.
	stored, err := s.repository.GetManagedComplex(ctx, tx, input.ID, input.Kind)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.update.get_error", "err", err, "id", input.ID)
		return nil, utils.InternalError("")
	}

	if stored.Kind() != input.Kind {
		return nil, utils.ValidationError("coverageType", "Coverage type cannot be changed.")
	}

	rows, err := s.repository.UpdateManagedComplex(ctx, tx, domain)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.update.repo_error", "err", err)
		return nil, utils.InternalError("")
	}

	if rows == 0 {
		return nil, utils.NotFoundError("complex")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.update.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true
	return domain, nil
}
