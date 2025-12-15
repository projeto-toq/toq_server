package photosessionservices

import (
	"fmt"
	"strings"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func validateTimeOffInput(input TimeOffInput) error {
	if input.PhotographerID == 0 {
		return utils.ValidationError("photographerId", "photographerId must be greater than zero")
	}
	if input.StartDate.IsZero() {
		return utils.ValidationError("startDate", "startDate is required")
	}
	if input.EndDate.IsZero() {
		return utils.ValidationError("endDate", "endDate is required")
	}
	if input.Location == nil {
		return utils.ValidationError("timezone", "timezone is required")
	}
	start := utils.ConvertToLocation(input.StartDate, input.Location)
	end := utils.ConvertToLocation(input.EndDate, input.Location)
	if end.Before(start) {
		return utils.ValidationError("endDate", "endDate must be greater than or equal to startDate")
	}
	if input.Reason != nil {
		reason := strings.TrimSpace(*input.Reason)
		if len(reason) > maxTimeOffReasonLength {
			return utils.ValidationError("reason", fmt.Sprintf("reason must be at most %d characters", maxTimeOffReasonLength))
		}
	}
	return nil
}

func validateEnsureAgendaInput(input EnsureAgendaInput) error {
	if input.PhotographerID == 0 {
		return utils.ValidationError("photographerId", "photographerId must be greater than zero")
	}
	if input.Timezone != "" {
		if _, err := resolveLocation(input.Timezone); err != nil {
			return err
		}
	}
	return nil
}

func validateListAgendaInput(input ListAgendaInput) error {
	if input.PhotographerID == 0 {
		return derrors.Validation("photographerId must be greater than zero", nil)
	}
	if input.StartDate.IsZero() {
		return derrors.Validation("startDate is required", nil)
	}
	if input.EndDate.IsZero() {
		return derrors.Validation("endDate is required", nil)
	}
	if input.EndDate.Before(input.StartDate) {
		return derrors.Validation("endDate must be after or equal to startDate", nil)
	}
	if input.Size < 0 {
		return derrors.Validation("size must be zero or greater", nil)
	}
	if input.Page < 0 {
		return derrors.Validation("page must be zero or greater", nil)
	}
	if input.SortField != "" {
		switch input.SortField {
		case AgendaSortFieldStartDate, AgendaSortFieldEndDate, AgendaSortFieldEntryType:
		default:
			return derrors.Validation("invalid sortField", map[string]any{"sortField": input.SortField})
		}
	}
	if input.SortOrder != "" {
		switch input.SortOrder {
		case AgendaSortOrderAsc, AgendaSortOrderDesc:
		default:
			return derrors.Validation("invalid sortOrder", map[string]any{"sortOrder": input.SortOrder})
		}
	}
	return nil
}

func resolveLocation(timezone string) (*time.Location, error) {
	if timezone == "" {
		timezone = defaultTimezone
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, utils.ValidationError("timezone", "Invalid timezone")
	}

	return loc, nil
}

func listingAllowsPhotoSession(status listingmodel.ListingStatus) bool {
	switch status {
	case listingmodel.StatusPendingPhotoScheduling,
		listingmodel.StatusPendingPhotoConfirmation,
		listingmodel.StatusPhotosScheduled:
		return true
	default:
		return false
	}
}
