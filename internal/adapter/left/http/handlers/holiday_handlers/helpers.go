package holidayhandlers

import (
	"strings"
	"time"

	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	defaultHolidayLimit = 20
)

func parseHolidayScope(value string) (holidaymodel.CalendarScope, error) {
	val := strings.ToUpper(strings.TrimSpace(value))
	switch holidaymodel.CalendarScope(val) {
	case holidaymodel.ScopeNational,
		holidaymodel.ScopeState,
		holidaymodel.ScopeCity:
		return holidaymodel.CalendarScope(val), nil
	case "":
		return "", nil
	default:
		return "", utils.ValidationError("scope", "scope must be NATIONAL, STATE or CITY")
	}
}

func parseHolidayDate(field, value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Time{}, utils.ValidationError(field, field+" is required and must be RFC3339 timestamp")
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, utils.ValidationError(field, field+" must be a valid RFC3339 timestamp")
	}
	return parsed, nil
}

func parseOptionalHolidayDate(field, value string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, utils.ValidationError(field, field+" must be a valid RFC3339 timestamp")
	}
	return &parsed, nil
}

func sanitizeHolidayPagination(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = defaultHolidayLimit
	}
	return page, limit
}
