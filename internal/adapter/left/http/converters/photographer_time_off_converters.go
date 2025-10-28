package converters

import (
	"math"
	"time"

	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	photosessionmodel "github.com/projeto-toq/toq_server/internal/core/model/photo_session_model"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

// PhotographerTimeOffToDTO converts a time-off domain entity into a response DTO.
func PhotographerTimeOffToDTO(timeOff photosessionmodel.PhotographerTimeOffInterface, timezone string) dto.PhotographerTimeOffResponse {
	if timeOff == nil {
		return dto.PhotographerTimeOffResponse{}
	}

	resp := dto.PhotographerTimeOffResponse{
		ID:        timeOff.ID(),
		StartDate: formatTimeOff(timeOff.StartDate()),
		EndDate:   formatTimeOff(timeOff.EndDate()),
		Timezone:  timezone,
	}

	if reason := timeOff.Reason(); reason != nil {
		resp.Reason = reason
	}

	return resp
}

// ListTimeOffOutputToDTO converts the service list result into the HTTP response payload.
func ListTimeOffOutputToDTO(output photosessionservices.ListTimeOffOutput) dto.ListPhotographerTimeOffResponse {
	items := make([]dto.PhotographerTimeOffResponse, 0, len(output.TimeOffs))
	for _, entry := range output.TimeOffs {
		items = append(items, PhotographerTimeOffToDTO(entry, output.Timezone))
	}

	return dto.ListPhotographerTimeOffResponse{
		TimeOffs:   items,
		Pagination: buildTimeOffPagination(output.Page, output.Size, output.Total),
		Timezone:   output.Timezone,
	}
}

// TimeOffResultToDTO converts a single time-off service result into DTO.
func TimeOffResultToDTO(result photosessionservices.TimeOffDetailResult) dto.PhotographerTimeOffResponse {
	return PhotographerTimeOffToDTO(result.TimeOff, result.Timezone)
}

func buildTimeOffPagination(page, size int, total int64) dto.PaginationResponse {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	return dto.PaginationResponse{
		Page:       page,
		Limit:      size,
		Total:      total,
		TotalPages: calculateTimeOffTotalPages(total, size),
	}
}

func calculateTimeOffTotalPages(total int64, size int) int {
	if size <= 0 {
		return 1
	}
	if total == 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(size)))
}

func formatTimeOff(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339)
}
