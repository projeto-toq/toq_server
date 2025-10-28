package dto

// ListAgendaQuery captures query parameters for listing the photographer agenda.
type ListAgendaQuery struct {
	StartDate          string `form:"startDate" binding:"required" example:"2023-10-01T00:00:00Z"`
	EndDate            string `form:"endDate" binding:"required" example:"2023-10-31T23:59:59Z"`
	Page               int    `form:"page" binding:"omitempty,min=1" example:"1"`
	Size               int    `form:"size" binding:"omitempty,min=1" example:"20"`
	Timezone           string `form:"timezone" binding:"omitempty" example:"America/Sao_Paulo"`
	HolidayCalendarIDs string `form:"holidayCalendarIds" binding:"omitempty" example:"1,2,3"`
}

// UpdateSessionStatusRequest defines the payload for updating a session's status.
type UpdateSessionStatusRequest struct {
	SessionID uint64 `json:"sessionId" binding:"required" example:"12345"`
	Status    string `json:"status" binding:"required" example:"ACCEPTED"`
}

// CreateTimeOffRequest represents the payload to block a photographer agenda.
type CreateTimeOffRequest struct {
	StartDate         string  `json:"startDate" binding:"required" example:"2023-11-01T09:00:00-03:00"`
	EndDate           string  `json:"endDate" binding:"required" example:"2023-11-01T18:00:00-03:00"`
	Reason            *string `json:"reason,omitempty" example:"Attending workshop"`
	Timezone          string  `json:"timezone" binding:"required" example:"America/Sao_Paulo"`
	HolidayCalendarID *uint64 `json:"holidayCalendarId,omitempty" example:"1"`
	HorizonMonths     int     `json:"horizonMonths" binding:"required" example:"2"`
	WorkdayStartHour  int     `json:"workdayStartHour" binding:"required" example:"9"`
	WorkdayEndHour    int     `json:"workdayEndHour" binding:"required" example:"18"`
}

// DeleteTimeOffRequest represents the payload to unblock a photographer agenda.
type DeleteTimeOffRequest struct {
	TimeOffID         uint64  `json:"timeOffId" binding:"required" example:"42"`
	Timezone          string  `json:"timezone" binding:"required" example:"America/Sao_Paulo"`
	HolidayCalendarID *uint64 `json:"holidayCalendarId,omitempty" example:"1"`
	HorizonMonths     int     `json:"horizonMonths" binding:"required" example:"2"`
	WorkdayStartHour  int     `json:"workdayStartHour" binding:"required" example:"9"`
	WorkdayEndHour    int     `json:"workdayEndHour" binding:"required" example:"18"`
}

// ListTimeOffQuery captures filters to list photographer time-offs.
type ListTimeOffQuery struct {
	RangeFrom string `form:"rangeFrom" binding:"required" example:"2025-07-01T00:00:00Z"`
	RangeTo   string `form:"rangeTo" binding:"required" example:"2025-07-31T23:59:59Z"`
	Page      int    `form:"page" binding:"omitempty,min=1" example:"1"`
	Size      int    `form:"size" binding:"omitempty,min=1" example:"20"`
	Timezone  string `form:"timezone" binding:"omitempty" example:"America/Sao_Paulo"`
}

// UpdateTimeOffRequest represents the payload to update a photographer time-off.
type UpdateTimeOffRequest struct {
	TimeOffID         uint64  `json:"timeOffId" binding:"required" example:"42"`
	StartDate         string  `json:"startDate" binding:"required" example:"2025-07-05T09:00:00-03:00"`
	EndDate           string  `json:"endDate" binding:"required" example:"2025-07-05T18:00:00-03:00"`
	Reason            *string `json:"reason,omitempty" example:"Atualizacao de agenda"`
	Timezone          string  `json:"timezone" binding:"required" example:"America/Sao_Paulo"`
	HolidayCalendarID *uint64 `json:"holidayCalendarId,omitempty" example:"1"`
	HorizonMonths     int     `json:"horizonMonths" binding:"required" example:"3"`
	WorkdayStartHour  int     `json:"workdayStartHour" binding:"required" example:"8"`
	WorkdayEndHour    int     `json:"workdayEndHour" binding:"required" example:"19"`
}

// TimeOffDetailRequest carries identifiers to fetch a specific time-off.
type TimeOffDetailRequest struct {
	TimeOffID uint64 `json:"timeOffId" binding:"required" example:"42"`
	Timezone  string `json:"timezone" binding:"omitempty" example:"America/Sao_Paulo"`
}

// PhotographerTimeOffResponse represents a normalized time-off entry.
type PhotographerTimeOffResponse struct {
	ID        uint64  `json:"id"`
	StartDate string  `json:"startDate"`
	EndDate   string  `json:"endDate"`
	Reason    *string `json:"reason,omitempty"`
	Timezone  string  `json:"timezone"`
}

// ListPhotographerTimeOffResponse aggregates time-off entries for a photographer.
type ListPhotographerTimeOffResponse struct {
	TimeOffs   []PhotographerTimeOffResponse `json:"timeOffs"`
	Pagination PaginationResponse            `json:"pagination"`
	Timezone   string                        `json:"timezone"`
}
