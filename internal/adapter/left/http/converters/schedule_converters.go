package converters

import (
	"fmt"
	"math"
	"strings"
	"time"

	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	scheduleservices "github.com/projeto-toq/toq_server/internal/core/service/schedule_service"
)

// ScheduleOwnerSummaryToDTO converts the owner summary domain result into a response DTO.
func ScheduleOwnerSummaryToDTO(result schedulemodel.OwnerSummaryResult, page, limit int) dto.OwnerAgendaSummaryResponse {
	items := make([]dto.OwnerAgendaSummaryItemResponse, 0, len(result.Items))
	for _, item := range result.Items {
		summaryEntries := make([]dto.OwnerAgendaSummaryEntryResponse, 0, len(item.Entries))
		for _, entry := range item.Entries {
			summaryEntries = append(summaryEntries, dto.OwnerAgendaSummaryEntryResponse{
				EntryType: string(entry.EntryType),
				StartsAt:  formatScheduleTime(entry.StartsAt),
				EndsAt:    formatScheduleTime(entry.EndsAt),
				Blocking:  entry.Blocking,
			})
		}
		items = append(items, dto.OwnerAgendaSummaryItemResponse{
			ListingID: item.ListingID,
			Entries:   summaryEntries,
		})
	}

	return dto.OwnerAgendaSummaryResponse{
		Items:      items,
		Pagination: buildSchedulePagination(page, limit, result.Total),
	}
}

// ScheduleEntriesToDTO converts agenda timeline items into a response DTO.
func ScheduleEntriesToDTO(result schedulemodel.AgendaDetailResult, page, limit int) dto.ListingAgendaDetailResponse {
	detailEntries := make([]dto.ScheduleEntryResponse, 0, len(result.Items))
	for _, item := range result.Items {
		detailEntries = append(detailEntries, timelineItemToDTO(item, result.Timezone))
	}

	return dto.ListingAgendaDetailResponse{
		Entries:    detailEntries,
		Pagination: buildSchedulePagination(page, limit, result.Total),
		Timezone:   result.Timezone,
	}
}

// ScheduleBlockEntriesToDTO converts blocking agenda entries into a response DTO.
func ScheduleBlockEntriesToDTO(result schedulemodel.BlockEntriesResult, page, limit int) dto.ListingBlockEntriesResponse {
	entries := make([]dto.ScheduleEntryResponse, 0, len(result.Items))
	for _, entry := range result.Items {
		if entry == nil {
			continue
		}
		entries = append(entries, ScheduleEntryToDTO(entry, result.Timezone))
	}

	return dto.ListingBlockEntriesResponse{
		Entries:    entries,
		Pagination: buildSchedulePagination(page, limit, result.Total),
		Timezone:   result.Timezone,
	}
}

// ScheduleEntryToDTO converts a domain agenda entry into a DTO representation.
func ScheduleEntryToDTO(entry schedulemodel.AgendaEntryInterface, fallbackTimezone string) dto.ScheduleEntryResponse {
	if entry == nil {
		return dto.ScheduleEntryResponse{}
	}

	response := dto.ScheduleEntryResponse{
		ID:         entry.ID(),
		SourceType: string(schedulemodel.TimelineSourceEntry),
		Recurring:  false,
		EntryType:  string(entry.EntryType()),
		StartsAt:   formatScheduleTime(entry.StartsAt()),
		EndsAt:     formatScheduleTime(entry.EndsAt()),
		Blocking:   entry.Blocking(),
	}

	if loc := entry.StartsAt().Location(); loc != nil {
		response.Timezone = loc.String()
	} else {
		response.Timezone = fallbackTimezone
	}

	response.Weekday = formatWeekday(entry.StartsAt().Weekday())

	if reason, ok := entry.Reason(); ok {
		response.Reason = reason
	}

	if visitID, ok := entry.VisitID(); ok {
		response.VisitID = visitID
	}

	if photoID, ok := entry.PhotoBookingID(); ok {
		response.PhotoBookingID = photoID
	}

	return response
}

func timelineItemToDTO(item schedulemodel.AgendaTimelineItem, fallbackTimezone string) dto.ScheduleEntryResponse {
	if item.Entry != nil {
		return ScheduleEntryToDTO(item.Entry, fallbackTimezone)
	}

	response := dto.ScheduleEntryResponse{
		SourceType: string(item.Source),
		Recurring:  item.Recurring,
		Weekday:    formatWeekday(item.Weekday),
		StartsAt:   formatScheduleTime(item.StartsAt),
		EndsAt:     formatScheduleTime(item.EndsAt),
		Blocking:   item.Blocking,
		Timezone:   fallbackTimezone,
	}

	if item.Rule != nil {
		response.RuleID = item.Rule.ID()
	}

	return response
}

func formatWeekday(day time.Weekday) string {
	if day < 0 || day > 6 {
		return ""
	}
	return strings.ToUpper(day.String())
}

// ScheduleRulesMutationToDTO converts a rule mutation result into a response payload.
func ScheduleRulesMutationToDTO(result scheduleservices.RuleMutationResult) dto.ScheduleRulesResponse {
	return dto.ScheduleRulesResponse{
		ListingID: result.ListingID,
		Timezone:  result.Timezone,
		Rules:     scheduleRulesToDTO(result.Rules),
	}
}

// ScheduleRuleListToDTO converts a rule list domain result into a response payload.
func ScheduleRuleListToDTO(result schedulemodel.RuleListResult) dto.ScheduleRulesResponse {
	return dto.ScheduleRulesResponse{
		ListingID: result.ListingID,
		Timezone:  result.Timezone,
		Rules:     scheduleRulesToDTO(result.Rules),
	}
}

// ScheduleRuleToDTO converts a single domain rule into a response representation.
func ScheduleRuleToDTO(rule schedulemodel.AgendaRuleInterface) dto.ScheduleRuleResponse {
	if rule == nil {
		return dto.ScheduleRuleResponse{}
	}
	return dto.ScheduleRuleResponse{
		RuleID:    rule.ID(),
		Weekday:   formatWeekday(rule.DayOfWeek()),
		StartTime: formatMinutesAsTime(rule.StartMinutes()),
		EndTime:   formatMinutesAsTime(rule.EndMinutes()),
		Active:    rule.IsActive(),
	}
}

func scheduleRulesToDTO(rules []schedulemodel.AgendaRuleInterface) []dto.ScheduleRuleResponse {
	responses := make([]dto.ScheduleRuleResponse, 0, len(rules))
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		responses = append(responses, ScheduleRuleToDTO(rule))
	}
	return responses
}

func formatMinutesAsTime(minutes uint16) string {
	hour := int(minutes) / 60
	minute := int(minutes) % 60
	return fmt.Sprintf("%02d:%02d", hour, minute)
}

// ScheduleAvailabilityToDTO converts the availability result into a response DTO.
func ScheduleAvailabilityToDTO(result scheduleservices.AvailabilityResult, page, limit int) dto.ScheduleAvailabilityResponse {
	slots := make([]dto.ScheduleAvailabilitySlotResponse, 0, len(result.Slots))
	for _, slot := range result.Slots {
		slots = append(slots, dto.ScheduleAvailabilitySlotResponse{
			StartsAt: formatScheduleTime(slot.StartsAt),
			EndsAt:   formatScheduleTime(slot.EndsAt),
		})
	}

	return dto.ScheduleAvailabilityResponse{
		Slots:      slots,
		Pagination: buildSchedulePagination(page, limit, int64(result.Total)),
		Timezone:   result.Timezone,
	}
}

func buildSchedulePagination(page, limit int, total int64) dto.PaginationResponse {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	return dto.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: calculateTotalPages(total, limit),
	}
}

func calculateTotalPages(total int64, limit int) int {
	if limit <= 0 {
		return 1
	}
	if total == 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}

func formatScheduleTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339)
}
