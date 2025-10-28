package dto

// HolidayCalendarRequest captures data to create or update a holiday calendar.
type HolidayCalendarRequest struct {
	Name     string `json:"name" binding:"required"`
	Scope    string `json:"scope" binding:"required"`
	State    string `json:"state,omitempty"`
	CityIBGE string `json:"cityIbge,omitempty"`
	IsActive bool   `json:"isActive"`
	Timezone string `json:"timezone" binding:"required" example:"America/Sao_Paulo"`
}

// HolidayCalendarCreateRequest extends the calendar request for creation.
type HolidayCalendarCreateRequest struct {
	HolidayCalendarRequest
}

// HolidayCalendarUpdateRequest extends the calendar request for updates.
type HolidayCalendarUpdateRequest struct {
	ID uint64 `json:"id" binding:"required"`
	HolidayCalendarRequest
}

// HolidayCalendarDetailRequest holds the identifier to fetch a calendar.
type HolidayCalendarDetailRequest struct {
	ID uint64 `json:"id" binding:"required"`
}

// HolidayCalendarsListRequest contains query parameters to list calendars.
type HolidayCalendarsListRequest struct {
	Scope      string `form:"scope"`
	State      string `form:"state"`
	CityIBGE   string `form:"cityIbge"`
	SearchTerm string `form:"search"`
	OnlyActive *bool  `form:"onlyActive"`
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=20"`
}

// HolidayCalendarResponse represents a holiday calendar in responses.
type HolidayCalendarResponse struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Scope    string `json:"scope"`
	State    string `json:"state,omitempty"`
	CityIBGE string `json:"cityIbge,omitempty"`
	IsActive bool   `json:"isActive"`
	Timezone string `json:"timezone"`
}

// HolidayCalendarsListResponse aggregates calendars with pagination metadata.
type HolidayCalendarsListResponse struct {
	Calendars  []HolidayCalendarResponse `json:"calendars"`
	Pagination PaginationResponse        `json:"pagination"`
}

// HolidayCalendarDateCreateRequest captures data to create a holiday date entry.
type HolidayCalendarDateCreateRequest struct {
	CalendarID  uint64 `json:"calendarId" binding:"required"`
	HolidayDate string `json:"holidayDate" binding:"required"`
	Label       string `json:"label" binding:"required"`
	Recurrent   bool   `json:"recurrent"`
}

// HolidayCalendarDatesListRequest contains query parameters to list calendar dates.
type HolidayCalendarDatesListRequest struct {
	CalendarID uint64 `form:"calendarId" binding:"required"`
	From       string `form:"from"`
	To         string `form:"to"`
	Timezone   string `form:"timezone"`
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=20"`
}

// HolidayCalendarDateResponse represents a holiday date in responses.
type HolidayCalendarDateResponse struct {
	ID          uint64 `json:"id"`
	CalendarID  uint64 `json:"calendarId"`
	HolidayDate string `json:"holidayDate"`
	Label       string `json:"label"`
	Recurrent   bool   `json:"recurrent"`
	Timezone    string `json:"timezone"`
}

// HolidayCalendarDatesListResponse aggregates holiday dates with pagination metadata.
type HolidayCalendarDatesListResponse struct {
	Dates      []HolidayCalendarDateResponse `json:"dates"`
	Pagination PaginationResponse            `json:"pagination"`
}

// HolidayCalendarDateDeleteRequest carries the identifier to remove a calendar date.
type HolidayCalendarDateDeleteRequest struct {
	ID uint64 `json:"id" binding:"required"`
}
