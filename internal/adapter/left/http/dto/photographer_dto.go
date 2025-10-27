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
