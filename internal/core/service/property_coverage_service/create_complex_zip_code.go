package propertycoverageservice

import (
	"context"
	"database/sql"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateComplexZipCode registers a zip code entry for a horizontal complex.
func (s *propertyCoverageService) CreateComplexZipCode(ctx context.Context, input CreateComplexZipCodeInput) (propertycoveragemodel.HorizontalComplexZipCodeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	if err := ensurePositiveID("horizontalComplexId", input.HorizontalComplexID); err != nil {
		return nil, err
	}

	normalizedZip, err := normalizeAndValidateZip(input.ZipCode)
	if err != nil {
		return nil, err
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.zip.create.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.zip.create.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	if _, err := s.repository.GetManagedComplex(ctx, tx, input.HorizontalComplexID, propertycoveragemodel.CoverageKindHorizontal); err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("complex")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.zip.create.parent_error", "err", err, "complex_id", input.HorizontalComplexID)
		return nil, utils.InternalError("")
	}

	zipEntity := propertycoveragemodel.NewHorizontalComplexZipCode()
	zipEntity.SetHorizontalComplexID(input.HorizontalComplexID)
	zipEntity.SetZipCode(normalizedZip)

	id, err := s.repository.CreateHorizontalComplexZipCode(ctx, tx, zipEntity)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.zip.create.repo_error", "err", err)
		return nil, utils.InternalError("")
	}

	zipEntity.SetID(id)

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.zip.create.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true
	return zipEntity, nil
}
