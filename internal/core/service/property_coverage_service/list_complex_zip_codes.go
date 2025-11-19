package propertycoverageservice

import (
	"context"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	propertycoveragerepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/property_coverage_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

// ListComplexZipCodes lists zip code entries for a horizontal complex.
func (s *propertyCoverageService) ListComplexZipCodes(ctx context.Context, input ListComplexZipCodesInput) ([]propertycoveragemodel.HorizontalComplexZipCodeInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	if err := ensurePositiveID("horizontalComplexId", input.HorizontalComplexID); err != nil {
		return nil, err
	}

	zipFilter := sanitizeString(input.ZipCode)
	if zipFilter != "" {
		normalized, normErr := validators.NormalizeCEP(zipFilter)
		if normErr != nil {
			return nil, utils.ValidationError("zipCode", "Zip code must contain exactly 8 digits without separators.")
		}
		zipFilter = normalized
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	page, limit := sanitizePagination(input.Page, input.Limit)
	offset := (page - 1) * limit

	params := propertycoveragerepository.ListHorizontalComplexZipCodesParams{
		HorizontalComplexID: input.HorizontalComplexID,
		ZipCode:             zipFilter,
		Limit:               limit,
		Offset:              offset,
	}

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.zip.list.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.zip.list.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	zips, err := s.repository.ListHorizontalComplexZipCodes(ctx, tx, params)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.zip.list.repo_error", "err", err)
		return nil, utils.InternalError("")
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.zip.list.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true
	return zips, nil
}
