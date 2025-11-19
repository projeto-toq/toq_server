package propertycoverageservice

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

// ResolvePropertyTypes determines the allowed property types for the provided zip code/number.
func (s *propertyCoverageService) ResolvePropertyTypes(ctx context.Context, input propertycoveragemodel.ResolvePropertyTypesInput) (output propertycoveragemodel.ResolvePropertyTypesOutput, err error) {
	ctx, spanEnd, spanErr := utils.GenerateTracer(ctx)
	if spanErr != nil {
		return output, utils.InternalError("Failed to initialize tracing")
	}
	defer spanEnd()

	zipCode, number, normErr := normalizeCoverageInput(input)
	if normErr != nil {
		return output, normErr
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := s.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("property_coverage.resolve.tx_start_error", "err", err)
		return output, utils.InternalError("Failed to start transaction")
	}

	defer func() {
		if err != nil {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("property_coverage.resolve.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	var coverage propertycoveragemodel.CoverageInterface
	var lookupErr error

	if number != "" {
		coverage, lookupErr = s.repository.GetVerticalCoverage(ctx, tx, zipCode, number)
		if lookupErr != nil {
			if !errors.Is(lookupErr, sql.ErrNoRows) {
				utils.SetSpanError(ctx, lookupErr)
				logger.Error("property_coverage.resolve.vertical_error", "zip_code", zipCode, "number", number, "err", lookupErr)
				return output, utils.InternalError("")
			}
		} else {
			goto finalize
		}
	}

	coverage, lookupErr = s.repository.GetHorizontalCoverage(ctx, tx, zipCode)
	if lookupErr != nil {
		if errors.Is(lookupErr, sql.ErrNoRows) {
			coverage, lookupErr = s.repository.GetNoComplexCoverage(ctx, tx, zipCode)
			if lookupErr != nil {
				if errors.Is(lookupErr, sql.ErrNoRows) {
					err = utils.NewHTTPErrorWithSource(http.StatusNotFound, "Area not covered yet for the provided zip code and number.")
					return output, err
				}
				utils.SetSpanError(ctx, lookupErr)
				logger.Error("property_coverage.resolve.no_complex_error", "zip_code", zipCode, "err", lookupErr)
				return output, utils.InternalError("")
			}
		} else {
			utils.SetSpanError(ctx, lookupErr)
			logger.Error("property_coverage.resolve.horizontal_error", "zip_code", zipCode, "err", lookupErr)
			return output, utils.InternalError("")
		}
	}

finalize:
	output = coverage.ToOutput()

	if commitErr := s.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("property_coverage.resolve.tx_commit_error", "err", commitErr)
		return output, utils.InternalError("Failed to commit transaction")
	}

	return output, nil
}

func normalizeCoverageInput(input propertycoveragemodel.ResolvePropertyTypesInput) (zipCode string, number string, err error) {
	zipCandidate := strings.TrimSpace(input.ZipCode)
	zipCode, err = validators.NormalizeCEP(zipCandidate)
	if err != nil {
		return "", "", utils.ValidationError("zipCode", "Zip code must contain exactly 8 digits without separators.")
	}

	number = sanitizeCoverageNumber(input.Number)
	return zipCode, number, nil
}

func sanitizeCoverageNumber(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	noSpaces := strings.ReplaceAll(trimmed, " ", "")
	return strings.ToUpper(noSpaces)
}
