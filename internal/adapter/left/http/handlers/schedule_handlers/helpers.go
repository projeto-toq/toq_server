package schedulehandlers

import (
	"strconv"
	"strings"
	"time"

	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	defaultScheduleLimit = 20
	minutesPerDay        = 24 * 60
)

var weekdayLookup = map[string]time.Weekday{
	"SUNDAY":    time.Sunday,
	"MONDAY":    time.Monday,
	"TUESDAY":   time.Tuesday,
	"WEDNESDAY": time.Wednesday,
	"THURSDAY":  time.Thursday,
	"FRIDAY":    time.Friday,
	"SATURDAY":  time.Saturday,
}

func parseScheduleRange(rangeReq dto.ScheduleRangeRequest) (schedulemodel.ScheduleRange, error) {
	var result schedulemodel.ScheduleRange
	if rangeReq.From != "" {
		from, err := time.Parse(time.RFC3339, rangeReq.From)
		if err != nil {
			return schedulemodel.ScheduleRange{}, utils.ValidationError("range.from", "from must be a valid RFC3339 timestamp")
		}
		result.From = from
	}
	if rangeReq.To != "" {
		to, err := time.Parse(time.RFC3339, rangeReq.To)
		if err != nil {
			return schedulemodel.ScheduleRange{}, utils.ValidationError("range.to", "to must be a valid RFC3339 timestamp")
		}
		result.To = to
	}
	return result, nil
}

func sanitizeSchedulePagination(p dto.SchedulePaginationRequest) schedulemodel.PaginationConfig {
	page := p.Page
	limit := p.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = defaultScheduleLimit
	}
	return schedulemodel.PaginationConfig{Page: page, Limit: limit}
}

func parseScheduleEntryType(value string) (schedulemodel.EntryType, error) {
	typed := schedulemodel.EntryType(strings.ToUpper(strings.TrimSpace(value)))
	switch typed {
	case schedulemodel.EntryTypeBlock, schedulemodel.EntryTypeTemporaryBlock:
		return typed, nil
	default:
		return "", utils.ValidationError("entryType", "entryType must be BLOCK or TEMP_BLOCK")
	}
}

func parseScheduleTimestamp(field, value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Time{}, utils.ValidationError(field, field+" is required and must be RFC3339 timestamp")
	}
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, utils.ValidationError(field, field+" must be a valid RFC3339 timestamp")
	}
	return ts, nil
}

func parseScheduleWeekdays(values []string) ([]time.Weekday, error) {
	if len(values) == 0 {
		return nil, utils.ValidationError("weekDays", "weekDays must contain at least one value")
	}

	result := make([]time.Weekday, 0, len(values))
	seen := make(map[time.Weekday]struct{})

	for _, raw := range values {
		key := strings.ToUpper(strings.TrimSpace(raw))
		day, ok := weekdayLookup[key]
		if !ok {
			return nil, utils.ValidationError("weekDays", "weekDays must contain valid weekday names (e.g. MONDAY)")
		}
		if _, exists := seen[day]; exists {
			continue
		}
		seen[day] = struct{}{}
		result = append(result, day)
	}

	return result, nil
}

func parseSingleScheduleWeekday(values []string) (time.Weekday, error) {
	weekdays, err := parseScheduleWeekdays(values)
	if err != nil {
		return 0, err
	}
	if len(weekdays) != 1 {
		return 0, utils.ValidationError("weekDays", "weekDays must contain exactly one value for this operation")
	}
	return weekdays[0], nil
}

func parseScheduleRuleMinutes(field, value string) (uint16, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0, utils.ValidationError(field, field+" is required and must be formatted as HH:MM")
	}

	parts := strings.Split(trimmed, ":")
	if len(parts) != 2 {
		return 0, utils.ValidationError(field, field+" must be formatted as HH:MM")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, utils.ValidationError(field, field+" must be formatted as HH:MM")
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, utils.ValidationError(field, field+" must be formatted as HH:MM")
	}

	if hour < 0 || hour > 24 {
		return 0, utils.ValidationError(field, field+" hour must be between 00 and 24")
	}
	if minute < 0 || minute >= 60 {
		return 0, utils.ValidationError(field, field+" minutes must be between 00 and 59")
	}
	if hour == 24 && minute != 0 {
		return 0, utils.ValidationError(field, field+" must be 24:00 or earlier")
	}

	total := hour*60 + minute
	if total > minutesPerDay {
		return 0, utils.ValidationError(field, field+" must be 24:00 or earlier")
	}

	return uint16(total), nil
}

func schedulePaginationValues(config schedulemodel.PaginationConfig) (page, limit int) {
	page = config.Page
	limit = config.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = defaultScheduleLimit
	}
	return
}
