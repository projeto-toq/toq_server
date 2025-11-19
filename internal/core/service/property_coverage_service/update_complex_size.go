package propertycoverageservice

import (
	"context"
	"database/sql"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateComplexSize mutates a size row linked to a vertical complex.
func (s *propertyCoverageService) UpdateComplexSize(ctx context.Context, input UpdateComplexSizeInput) (propertycoveragemodel.VerticalComplexSizeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	if err := ensurePositiveID("id", input.ID); err != nil {
		return nil, err
	}

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
		logger.Error("property_coverage.size.update.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.size.update.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	if _, err := s.repository.GetManagedComplex(ctx, tx, input.VerticalComplexID, propertycoveragemodel.CoverageKindVertical); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.size.update.parent_error", "err", err, "complex_id", input.VerticalComplexID)
		return nil, utils.InternalError("")
	}

	existing, err := s.repository.GetVerticalComplexSize(ctx, tx, input.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("size")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.size.update.get_error", "err", err, "id", input.ID)
		return nil, utils.InternalError("")
	}

	size := propertycoveragemodel.NewVerticalComplexSize()
	size.SetID(existing.ID())
	size.SetVerticalComplexID(input.VerticalComplexID)
	size.SetSize(input.Size)
	size.SetDescription(sanitizeString(input.Description))

	rows, err := s.repository.UpdateVerticalComplexSize(ctx, tx, size)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.size.update.repo_error", "err", err)
		return nil, utils.InternalError("")
	}

	if rows == 0 {
		return nil, utils.NotFoundError("size")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.size.update.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true
	return size, nil
}
