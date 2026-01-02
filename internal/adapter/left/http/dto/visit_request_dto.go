package dto

// CreateVisitRequest represents the payload to request a new visit.
// Lead time and horizon are enforced per configuration (visits.min_hours_ahead, visits.max_days_ahead).
type CreateVisitRequest struct {
	ListingIdentityID int64  `json:"listingIdentityId" binding:"required" example:"123"`
	ScheduledStart    string `json:"scheduledStart" binding:"required" example:"2025-01-10T14:00:00Z"`
	ScheduledEnd      string `json:"scheduledEnd" binding:"required" example:"2025-01-10T14:30:00Z"`
	Notes             string `json:"notes,omitempty" binding:"max=2000" example:"Client prefers afternoon"`
	Source            string `json:"source,omitempty" binding:"omitempty,oneof=APP WEB ADMIN" example:"APP"`
}

// UpdateVisitStatusRequest centralizes visit status transitions.
type UpdateVisitStatusRequest struct {
	VisitID         int64  `json:"visitId" binding:"required" example:"456"`
	Action          string `json:"action" binding:"required,oneof=APPROVE REJECT CANCEL COMPLETE NO_SHOW" example:"APPROVE"`
	RejectionReason string `json:"rejectionReason,omitempty" binding:"max=2000" example:"Slot unavailable"`
	Notes           string `json:"notes,omitempty" binding:"max=2000" example:"Owner approved with constraints"`
}

// VisitListQuery captures query parameters for visit listings (RFC3339 range, pagination capped at 50).
type VisitListQuery struct {
	ListingIdentityID int64    `form:"listingIdentityId"`
	Statuses          []string `form:"status"`
	From              string   `form:"from"`
	To                string   `form:"to"`
	Page              int      `form:"page,default=1"`
	Limit             int      `form:"limit,default=20"`
}
