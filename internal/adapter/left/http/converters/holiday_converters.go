package converters

import (
	"math"
	"time"

	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	holidaymodel "github.com/projeto-toq/toq_server/internal/core/model/holiday_model"
)

// HolidayCalendarToDTO converts a calendar domain object into a DTO.
func HolidayCalendarToDTO(calendar holidaymodel.CalendarInterface) dto.HolidayCalendarResponse {
	if calendar == nil {
		return dto.HolidayCalendarResponse{}
	}

	response := dto.HolidayCalendarResponse{
		ID:       calendar.ID(),
		Name:     calendar.Name(),
		Scope:    string(calendar.Scope()),
		IsActive: calendar.IsActive(),
		Timezone: calendar.Timezone(),
	}

	if state, ok := calendar.State(); ok {
		response.State = state
	}
	if city, ok := calendar.CityIBGE(); ok {
		response.CityIBGE = city
	}

	return response
}

// HolidayCalendarsToListDTO converts a list of calendars into a paginated response.
func HolidayCalendarsToListDTO(calendars []holidaymodel.CalendarInterface, page, limit int, total int64) dto.HolidayCalendarsListResponse {
	items := make([]dto.HolidayCalendarResponse, 0, len(calendars))
	for _, calendar := range calendars {
		items = append(items, HolidayCalendarToDTO(calendar))
	}

	return dto.HolidayCalendarsListResponse{
		Calendars:  items,
		Pagination: dto.PaginationResponse{Page: page, Limit: limit, Total: total, TotalPages: calculateHolidayTotalPages(total, limit)},
	}
}

// HolidayCalendarDateToDTO converts a calendar date domain object into a DTO.
func HolidayCalendarDateToDTO(date holidaymodel.CalendarDateInterface) dto.HolidayCalendarDateResponse {
	if date == nil {
		return dto.HolidayCalendarDateResponse{}
	}

	return dto.HolidayCalendarDateResponse{
		ID:          date.ID(),
		CalendarID:  date.CalendarID(),
		HolidayDate: formatHolidayDate(date.HolidayDate()),
		Label:       date.Label(),
		Recurrent:   date.IsRecurrent(),
		Timezone:    date.Timezone(),
	}
}

// HolidayCalendarDatesToListDTO converts domain holiday dates into a paginated DTO response.
func HolidayCalendarDatesToListDTO(dates []holidaymodel.CalendarDateInterface, page, limit int, total int64) dto.HolidayCalendarDatesListResponse {
	items := make([]dto.HolidayCalendarDateResponse, 0, len(dates))
	for _, date := range dates {
		items = append(items, HolidayCalendarDateToDTO(date))
	}

	return dto.HolidayCalendarDatesListResponse{
		Dates:      items,
		Pagination: dto.PaginationResponse{Page: page, Limit: limit, Total: total, TotalPages: calculateHolidayTotalPages(total, limit)},
	}
}

func calculateHolidayTotalPages(total int64, limit int) int {
	if limit <= 0 {
		return 1
	}
	if total == 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}

func formatHolidayDate(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339)
}
