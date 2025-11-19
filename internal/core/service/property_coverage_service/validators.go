package propertycoverageservice

import (
	"fmt"
	"strings"

	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

func validateRequiredField(fieldName, value string) error {
	if strings.TrimSpace(value) == "" {
		return utils.ValidationError(fieldName, fmt.Sprintf("%s is required.", fieldName))
	}
	return nil
}

func validateSector(sector propertycoveragemodel.Sector) error {
	switch sector {
	case propertycoveragemodel.SectorResidential,
		propertycoveragemodel.SectorCommercial,
		propertycoveragemodel.SectorMixed:
		return nil
	default:
		return utils.ValidationError("sector", "Invalid sector provided.")
	}
}

func normalizeAndValidateZip(zip string) (string, error) {
	candidate := strings.TrimSpace(zip)
	if candidate == "" {
		return "", utils.ValidationError("zipCode", "Zip code is required.")
	}

	normalized, err := validators.NormalizeCEP(candidate)
	if err != nil {
		return "", utils.ValidationError("zipCode", "Zip code must contain exactly 8 digits without separators.")
	}

	return normalized, nil
}

func normalizeOptional(value string) string {
	return strings.TrimSpace(value)
}

func ensurePositiveID(field string, id int64) error {
	if id <= 0 {
		return utils.ValidationError(field, fmt.Sprintf("%s must be greater than zero.", field))
	}
	return nil
}

func ensurePositiveFloat(field string, value float64) error {
	if value <= 0 {
		return utils.ValidationError(field, fmt.Sprintf("%s must be greater than zero.", field))
	}
	return nil
}

func validateCoverageKind(kind propertycoveragemodel.CoverageKind) error {
	if !kind.Valid() {
		return utils.ValidationError("coverageType", "Invalid coverage type provided.")
	}
	return nil
}
