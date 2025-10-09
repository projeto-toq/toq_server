package utils

import (
	"net/http"
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// ValidateUserDates parses and validates date fields from a UserCreateRequest.
// Returns parsed bornAt (required) and an optional creciValidity (if provided).
// On validation error, returns a DomainError (422) with field-level details using the given prefix (e.g., "owner", "realtor", "agency").
func ValidateUserDates(payload dto.UserCreateRequest, prefix string) (time.Time, *time.Time, coreutils.DomainError) {
	var details []map[string]string

	// bornAt is required by binding; still validate format here for precise error message
	const layout = "2006-01-02"
	bornAt, err := time.Parse(layout, payload.BornAt)
	if err != nil {
		details = append(details, map[string]string{
			"field":   prefix + ".bornAt",
			"message": "Invalid date format. Expected YYYY-MM-DD.",
		})
	}

	if err == nil {
		now := time.Now().UTC()
		age := now.Year() - bornAt.Year()
		birthdayThisYear := time.Date(now.Year(), bornAt.Month(), bornAt.Day(), 0, 0, 0, 0, time.UTC)
		if now.Before(birthdayThisYear) {
			age--
		}
		if age < 18 {
			details = append(details, map[string]string{
				"field":   prefix + ".bornAt",
				"message": "User must be at least 18 years old.",
			})
		}
	}

	var creciPtr *time.Time
	if payload.CreciValidity != "" {
		if t, err := time.Parse(layout, payload.CreciValidity); err != nil {
			details = append(details, map[string]string{
				"field":   prefix + ".creciValidity",
				"message": "Invalid date format. Expected YYYY-MM-DD.",
			})
		} else {
			creciPtr = &t
		}
	}

	if len(details) > 0 {
		return time.Time{}, nil, coreutils.NewHTTPError(http.StatusUnprocessableEntity, "Validation failed", map[string]any{
			"errors": details,
		})
	}

	return bornAt, creciPtr, nil
}
