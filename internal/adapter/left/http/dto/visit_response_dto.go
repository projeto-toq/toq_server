package dto

// VisitResponse represents a visit resource.
type VisitResponse struct {
	ID                 int64              `json:"id" example:"456"`
	ListingIdentityID  int64              `json:"listingIdentityId" example:"123"`
	ListingVersion     uint8              `json:"listingVersion" example:"1"`
	RequesterUserID    int64              `json:"requesterUserId" example:"5"`
	OwnerUserID        int64              `json:"ownerUserId" example:"10"`
	ScheduledStart     string             `json:"scheduledStart" example:"2025-01-10T14:00:00Z"`
	ScheduledEnd       string             `json:"scheduledEnd" example:"2025-01-10T14:30:00Z"`
	Status             string             `json:"status" example:"PENDING"`
	Source             string             `json:"source,omitempty" example:"APP"`
	Notes              string             `json:"notes,omitempty"`
	RejectionReason    string             `json:"rejectionReason,omitempty"`
	FirstOwnerActionAt *string            `json:"firstOwnerActionAt,omitempty" example:"2025-01-10T14:05:00Z"`
	ListingSummary     *ListingSummaryDTO `json:"listing,omitempty"`
}

// ListingSummaryDTO carries the essential listing data returned with a visit detail response.
type ListingSummaryDTO struct {
	Title        string `json:"title" example:"Cobertura incrível em Moema"`
	Description  string `json:"description" example:"Apartamento amplo com três suítes e vista livre."`
	ZipCode      string `json:"zipCode" example:"04534011"`
	Street       string `json:"street" example:"Av. Ibirapuera"`
	Number       string `json:"number" example:"1234"`
	Complement   string `json:"complement,omitempty" example:"apto 82"`
	Neighborhood string `json:"neighborhood" example:"Moema"`
	City         string `json:"city" example:"São Paulo"`
	State        string `json:"state" example:"SP"`
}

// VisitListResponse wraps paginated visit results.
type VisitListResponse struct {
	Items      []VisitResponse    `json:"items"`
	Pagination PaginationResponse `json:"pagination"`
}
