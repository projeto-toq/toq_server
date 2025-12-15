package dto

// ListAgendaQuery captures query parameters for listing the photographer agenda.
type ListAgendaQuery struct {
	StartDate string  `form:"startDate" binding:"required" example:"2023-10-01T00:00:00Z"`
	EndDate   string  `form:"endDate" binding:"required" example:"2023-10-31T23:59:59Z"`
	Page      int     `form:"page" binding:"omitempty,min=1" example:"1"`
	Size      int     `form:"size" binding:"omitempty,min=1" example:"20"`
	EntryType *string `form:"entryType" binding:"omitempty,oneof=PHOTO_SESSION BLOCK TIME_OFF HOLIDAY" example:"PHOTO_SESSION"`
	SortField string  `form:"sortField" binding:"omitempty,oneof=startDate endDate entryType" example:"startDate"`
	SortOrder string  `form:"sortOrder" binding:"omitempty,oneof=asc desc" example:"asc"`
}

// UpdateSessionStatusRequest defines the payload for updating a session's status.
// Status can be ACCEPTED (photographer accepts the session), REJECTED (photographer declines),
// or DONE (photographer confirms the session was completed).
type UpdateSessionStatusRequest struct {
	PhotoSessionID uint64 `json:"photoSessionId" binding:"required" example:"12345"`
	Status         string `json:"status" binding:"required,oneof=ACCEPTED REJECTED DONE" example:"ACCEPTED"`
}

// CreateTimeOffRequest represents the payload to block a photographer agenda.
type CreateTimeOffRequest struct {
	StartDate string  `json:"startDate" binding:"required" example:"2023-11-01T09:00:00-03:00"`
	EndDate   string  `json:"endDate" binding:"required" example:"2023-11-01T18:00:00-03:00"`
	Reason    *string `json:"reason,omitempty" example:"Attending workshop"`
}

// DeleteTimeOffRequest represents the payload to unblock a photographer agenda.
type DeleteTimeOffRequest struct {
	TimeOffID uint64 `json:"timeOffId" binding:"required" example:"42"`
}

// ListTimeOffQuery captures filters to list photographer time-offs.
type ListTimeOffQuery struct {
	RangeFrom string `form:"rangeFrom" binding:"required" example:"2025-07-01T00:00:00Z"`
	RangeTo   string `form:"rangeTo" binding:"required" example:"2025-07-31T23:59:59Z"`
	Page      int    `form:"page" binding:"omitempty,min=1" example:"1"`
	Size      int    `form:"size" binding:"omitempty,min=1" example:"20"`
}

// UpdateTimeOffRequest represents the payload to update a photographer time-off.
type UpdateTimeOffRequest struct {
	TimeOffID uint64  `json:"timeOffId" binding:"required" example:"42"`
	StartDate string  `json:"startDate" binding:"required" example:"2025-07-05T09:00:00-03:00"`
	EndDate   string  `json:"endDate" binding:"required" example:"2025-07-05T18:00:00-03:00"`
	Reason    *string `json:"reason,omitempty" example:"Agenda adjustment"`
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

// PhotographerServiceAreaListQuery captures filters for listing service areas.
type PhotographerServiceAreaListQuery struct {
	Page int `form:"page" binding:"omitempty,min=1" example:"1"`
	Size int `form:"size" binding:"omitempty,min=1" example:"20"`
}

// PhotographerServiceAreaRequest represents the payload to create or update a service area.
type PhotographerServiceAreaRequest struct {
	City  string `json:"city" binding:"required" example:"São Paulo"`
	State string `json:"state" binding:"required" example:"SP"`
}

// PhotographerServiceAreaDetailRequest carries the identifier to retrieve a service area.
type PhotographerServiceAreaDetailRequest struct {
	ServiceAreaID uint64 `json:"serviceAreaId" binding:"required" example:"42"`
}

// PhotographerServiceAreaUpdateRequest represents the payload to update a service area.
type PhotographerServiceAreaUpdateRequest struct {
	ServiceAreaID uint64 `json:"serviceAreaId" binding:"required" example:"42"`
	City          string `json:"city" binding:"required" example:"São Paulo"`
	State         string `json:"state" binding:"required" example:"SP"`
}

// PhotographerServiceAreaDeleteRequest carries the identifier to delete a service area.
type PhotographerServiceAreaDeleteRequest struct {
	ServiceAreaID uint64 `json:"serviceAreaId" binding:"required" example:"42"`
}

// PhotographerServiceAreaResponse represents a service area entry owned by the photographer.
type PhotographerServiceAreaResponse struct {
	ID    uint64 `json:"id" example:"42"`
	City  string `json:"city" example:"São Paulo"`
	State string `json:"state" example:"SP"`
}

// PhotographerServiceAreaListResponse aggregates service areas with pagination metadata.
type PhotographerServiceAreaListResponse struct {
	ServiceAreas []PhotographerServiceAreaResponse `json:"serviceAreas"`
	Pagination   PaginationResponse                `json:"pagination"`
}
