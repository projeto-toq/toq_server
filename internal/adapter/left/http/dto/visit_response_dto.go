package dto

// VisitResponse represents a visit resource enriched with owner/realtor metadata for both list and detail endpoints.
type VisitResponse struct {
	ID                 int64              `json:"id" example:"456"`
	ListingIdentityID  int64              `json:"listingIdentityId" example:"123"`
	ListingVersion     uint8              `json:"listingVersion" example:"1"`
	RequesterUserID    int64              `json:"requesterUserId" example:"5"`
	OwnerUserID        int64              `json:"ownerUserId" example:"10"`
	ScheduledStart     string             `json:"scheduledStart" example:"2025-01-10T14:00:00Z"`
	ScheduledEnd       string             `json:"scheduledEnd" example:"2025-01-10T14:30:00Z"`
	Status             string             `json:"status" example:"PENDING"`
	LiveStatus         string             `json:"liveStatus,omitempty" example:"AO_VIVO"`
	Source             string             `json:"source,omitempty" example:"APP"`
	Notes              string             `json:"notes,omitempty"`
	RejectionReason    string             `json:"rejectionReason,omitempty"`
	FirstOwnerActionAt *string            `json:"firstOwnerActionAt,omitempty" example:"2025-01-10T14:05:00Z"`
	ListingSummary     *ListingSummaryDTO `json:"listing,omitempty"`
	Owner              VisitOwnerDTO      `json:"owner"`
	Realtor            VisitRealtorDTO    `json:"realtor"`
	Timeline           VisitTimelineDTO   `json:"timeline"`
}

// ListingSummaryDTO carries the essential listing data returned with a visit detail response.
type ListingSummaryDTO struct {
	ListingIdentityID int64                        `json:"listingIdentityId" example:"123"`
	Title             string                       `json:"title" example:"Cobertura incrível em Moema"`
	Description       string                       `json:"description" example:"Apartamento amplo com três suítes e vista livre."`
	ZipCode           string                       `json:"zipCode" example:"04534011"`
	Street            string                       `json:"street" example:"Av. Ibirapuera"`
	Number            string                       `json:"number" example:"1234"`
	Complement        string                       `json:"complement,omitempty" example:"apto 82"`
	Neighborhood      string                       `json:"neighborhood" example:"Moema"`
	City              string                       `json:"city" example:"São Paulo"`
	State             string                       `json:"state" example:"SP"`
	PropertyType      *ListingPropertyTypeResponse `json:"propertyType,omitempty"`
}

// VisitListResponse wraps paginated visit results.
type VisitListResponse struct {
	Items      []VisitResponse    `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}

// VisitOwnerDTO exposes owner-facing metadata for realtor listings.
type VisitOwnerDTO struct {
	UserID           int64   `json:"userId" example:"10"`
	FullName         string  `json:"fullName" example:"Maria Souza"`
	PhotoURL         string  `json:"photoUrl,omitempty" example:"https://signed.cdn/photos/owner.jpg"`
	MemberSince      string  `json:"memberSince" example:"2021-05-10T12:00:00Z"`
	MemberSinceDays  int     `json:"memberSinceDays" example:"980"`
	AvgResponseHours float64 `json:"avgResponseHours,omitempty" example:"4.5"`
}

// VisitRealtorDTO exposes realtor metadata for owner listings.
type VisitRealtorDTO struct {
	UserID          int64  `json:"userId" example:"5"`
	FullName        string `json:"fullName" example:"João Corretor"`
	PhotoURL        string `json:"photoUrl,omitempty" example:"https://signed.cdn/photos/realtor.jpg"`
	MemberSince     string `json:"memberSince" example:"2022-02-01T09:30:00Z"`
	MemberSinceDays int    `json:"memberSinceDays" example:"600"`
	VisitsPerformed int64  `json:"visitsPerformed" example:"37"`
}

// VisitTimelineDTO keeps important timestamps for the visit lifecycle.
type VisitTimelineDTO struct {
	CreatedAt   string  `json:"createdAt" example:"2025-01-05T12:00:00Z"`
	ReceivedAt  string  `json:"receivedAt" example:"2025-01-05T12:05:00Z"`
	RespondedAt *string `json:"respondedAt,omitempty" example:"2025-01-05T13:15:00Z"`
}
