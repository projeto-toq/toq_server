package propertycoverageservice

import (
	"context"
	"database/sql"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateComplexSize registers a size entry for a vertical complex.
func (s *propertyCoverageService) CreateComplexSize(ctx context.Context, input CreateComplexSizeInput) (propertycoveragemodel.VerticalComplexSizeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	if err := ensurePositiveID("verticalComplexId", input.VerticalComplexID); err != nil {
		return nil, err
	}

	if err := ensurePositiveFloat("size", input.Size); err != nil {
		return nil, err
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.size.create.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.size.create.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	if _, err := s.repository.GetManagedComplex(ctx, tx, input.VerticalComplexID, propertycoveragemodel.CoverageKindVertical); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.size.create.parent_error", "err", err, "complex_id", input.VerticalComplexID)
		return nil, utils.InternalError("")
	}

	size := propertycoveragemodel.NewVerticalComplexSize()
	size.SetVerticalComplexID(input.VerticalComplexID)
	size.SetSize(input.Size)
	size.SetDescription(sanitizeString(input.Description))

	id, err := s.repository.CreateVerticalComplexSize(ctx, tx, size)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.size.create.repo_error", "err", err)
		return nil, utils.InternalError("")
	}

	size.SetID(id)

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.size.create.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true
	return size, nil
}
