package converters

import (
	"math"
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

// ScheduleEntriesToDTO converts agenda detail entries into a response DTO.
func ScheduleEntriesToDTO(entries []schedulemodel.AgendaDetailItem, page, limit int, total int64) dto.ListingAgendaDetailResponse {
	detailEntries := make([]dto.ScheduleEntryResponse, 0, len(entries))
	for _, item := range entries {
		if item.Entry == nil {
			continue
		}
		detailEntries = append(detailEntries, ScheduleEntryToDTO(item.Entry))
	}

	return dto.ListingAgendaDetailResponse{
		Entries:    detailEntries,
		Pagination: buildSchedulePagination(page, limit, total),
	}
}

// ScheduleEntryToDTO converts a domain agenda entry into a DTO representation.
func ScheduleEntryToDTO(entry schedulemodel.AgendaEntryInterface) dto.ScheduleEntryResponse {
	if entry == nil {
		return dto.ScheduleEntryResponse{}
	}

	response := dto.ScheduleEntryResponse{
		ID:        entry.ID(),
		EntryType: string(entry.EntryType()),
		StartsAt:  formatScheduleTime(entry.StartsAt()),
		EndsAt:    formatScheduleTime(entry.EndsAt()),
		Blocking:  entry.Blocking(),
	}

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
