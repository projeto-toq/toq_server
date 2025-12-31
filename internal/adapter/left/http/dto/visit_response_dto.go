package dto

// VisitResponse represents a visit resource.
type VisitResponse struct {
	ID                 int64   `json:"id" example:"456"`
	ListingIdentityID  int64   `json:"listingIdentityId" example:"123"`
	ListingVersion     uint8   `json:"listingVersion" example:"1"`
	RequesterUserID    int64   `json:"requesterUserId" example:"5"`
	OwnerUserID        int64   `json:"ownerUserId" example:"10"`
	ScheduledStart     string  `json:"scheduledStart" example:"2025-01-10T14:00:00Z"`
	ScheduledEnd       string  `json:"scheduledEnd" example:"2025-01-10T14:30:00Z"`
	DurationMinutes    int64   `json:"durationMinutes" example:"30"`
	Status             string  `json:"status" example:"PENDING"`
	Type               string  `json:"type" example:"WITH_CLIENT"`
	Source             string  `json:"source,omitempty" example:"APP"`
	RealtorNotes       string  `json:"realtorNotes,omitempty"`
	OwnerNotes         string  `json:"ownerNotes,omitempty"`
	RejectionReason    string  `json:"rejectionReason,omitempty"`
	CancelReason       string  `json:"cancelReason,omitempty"`
	FirstOwnerActionAt *string `json:"firstOwnerActionAt,omitempty" example:"2025-01-10T14:05:00Z"`
	CreatedAt          string  `json:"createdAt" example:"2025-01-09T12:00:00Z"`
	UpdatedAt          string  `json:"updatedAt" example:"2025-01-09T12:15:00Z"`
}

// VisitListResponse wraps paginated visit results.
type VisitListResponse struct {
	Items      []VisitResponse    `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}
