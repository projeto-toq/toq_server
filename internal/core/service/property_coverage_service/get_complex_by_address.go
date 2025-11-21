package propertycoverageservice

import (
	"context"
	"database/sql"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetComplexByAddress retrieves a managed complex by its address (ZipCode + Number).
// It searches for Vertical complexes first, then Horizontal.
// If neither is found, it checks for Standalone coverage (No Complex).
// Returns NotFoundError only if the area is completely uncovered.
func (s *propertyCoverageService) GetComplexByAddress(ctx context.Context, input GetComplexByAddressInput) (propertycoveragemodel.ManagedComplexInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	normalizedZip, err := normalizeAndValidateZip(input.ZipCode)
	if err != nil {
		return nil, err
	}

	number := sanitizeCoverageNumber(input.Number)

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("property_coverage.get_by_address.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	success := false
	defer func() {
		if !success {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.get_by_address.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	var complex propertycoveragemodel.ManagedComplexInterface

	// 1. Try Vertical Complex (requires Number)
	if number != "" {
		complex, err = s.repository.GetVerticalComplexByZipNumber(ctx, tx, normalizedZip, number)
		if err != nil && err != sql.ErrNoRows {
			utils.SetSpanError(ctx, err)
			logger.Error("property_coverage.get_by_address.get_vertical_error", "err", err, "zip_code", normalizedZip, "number", number)
			return nil, utils.InternalError("")
		}
	}

	// 2. If not found, try Horizontal Complex (ZipCode only)
	if complex == nil {
		complex, err = s.repository.GetHorizontalComplexByZip(ctx, tx, normalizedZip)
		if err != nil && err != sql.ErrNoRows {
			utils.SetSpanError(ctx, err)
			logger.Error("property_coverage.get_by_address.get_horizontal_error", "err", err, "zip_code", normalizedZip)
			return nil, utils.InternalError("")
		}
	}

	// 3. If still not found, try Standalone/NoComplex coverage
	if complex == nil {
		noComplexCoverage, err := s.repository.GetNoComplexCoverage(ctx, tx, normalizedZip)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, utils.NotFoundError("complex or coverage area")
			}
			utils.SetSpanError(ctx, err)
			logger.Error("property_coverage.get_by_address.get_nocomplex_error", "err", err, "zip_code", normalizedZip)
			return nil, utils.InternalError("")
		}

		// Create a virtual complex for standalone coverage
		complex = propertycoveragemodel.NewManagedComplex()
		complex.SetKind(propertycoveragemodel.CoverageKindStandalone)
		complex.SetPropertyTypes(noComplexCoverage.PropertyTypes())
		complex.SetZipCode(normalizedZip)
		complex.SetNumber(number)
		// ID remains 0, Name remains empty to indicate no specific complex
	}

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.get_by_address.tx_commit_error", "err", commitErr)
		return nil, utils.InternalError("")
	}

	success = true
	return complex, nil
}
