package schedulehandlers

import (
	"strings"
	"time"

	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	defaultScheduleLimit = 20
)

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
