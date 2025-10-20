package complexservices

import (
	"fmt"
	"strings"

	complexmodel "github.com/projeto-toq/toq_server/internal/core/model/complex_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

func validateRequiredField(fieldName, value string) error {
	if strings.TrimSpace(value) == "" {
		return utils.ValidationError(fieldName, fmt.Sprintf("%s is required.", fieldName))
	}
	return nil
}

func validateSector(sector complexmodel.Sector) error {
	switch sector {
	case complexmodel.SectorResidencial, complexmodel.SectorCommercial, complexmodel.SectorBoth:
		return nil
	default:
		return utils.ValidationError("sector", "Invalid sector provided.")
	}
}

func normalizeAndValidateZip(zip string) (string, error) {
	zipCandidate := strings.TrimSpace(zip)
	if zipCandidate == "" {
		return "", utils.ValidationError("zipCode", "Zip code is required.")
	}
	normalized, err := validators.NormalizeCEP(zipCandidate)
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

func sanitizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	return page, limit
}
