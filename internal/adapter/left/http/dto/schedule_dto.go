package dto

// ScheduleRangeRequest represents a time window filter.
type ScheduleRangeRequest struct {
	From string `json:"from,omitempty" example:"2025-01-01T00:00:00Z"`
	To   string `json:"to,omitempty" example:"2025-01-07T23:59:59Z"`
}

// SchedulePaginationRequest contains pagination parameters for schedule endpoints.
type SchedulePaginationRequest struct {
	Page  int `json:"page,omitempty" example:"1"`
	Limit int `json:"limit,omitempty" example:"20"`
}

// OwnerAgendaSummaryQuery captures query string parameters for owner agenda summary.
type OwnerAgendaSummaryQuery struct {
	ListingIDs []int64 `form:"listingIds"`
	RangeFrom  string  `form:"rangeFrom"`
	RangeTo    string  `form:"rangeTo"`
	Page       int     `form:"page"`
	Limit      int     `form:"limit"`
}

// OwnerAgendaSummaryEntryResponse describes a normalized agenda entry in the summary response.
type OwnerAgendaSummaryEntryResponse struct {
	EntryType string `json:"entryType"`
	StartsAt  string `json:"startsAt"`
	EndsAt    string `json:"endsAt"`
	Blocking  bool   `json:"blocking"`
}

// OwnerAgendaSummaryItemResponse groups summary entries for a specific listing.
type OwnerAgendaSummaryItemResponse struct {
	ListingID int64                             `json:"listingId"`
	Entries   []OwnerAgendaSummaryEntryResponse `json:"entries"`
}

// OwnerAgendaSummaryResponse aggregates the consolidated agenda view for owners.
type OwnerAgendaSummaryResponse struct {
	Items      []OwnerAgendaSummaryItemResponse `json:"items"`
	Pagination PaginationResponse               `json:"pagination"`
}

// ListingAgendaDetailQuery represents query parameters to list agenda entries of a specific listing.
type ListingAgendaDetailQuery struct {
	ListingID int64  `form:"listingId" binding:"required"`
	RangeFrom string `form:"rangeFrom"`
	RangeTo   string `form:"rangeTo"`
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
}

// ScheduleEntryResponse exposes detailed information about a single agenda entry.
type ScheduleEntryResponse struct {
	ID             uint64 `json:"id"`
	EntryType      string `json:"entryType"`
	StartsAt       string `json:"startsAt"`
	EndsAt         string `json:"endsAt"`
	Blocking       bool   `json:"blocking"`
	Reason         string `json:"reason,omitempty"`
	VisitID        uint64 `json:"visitId,omitempty"`
	PhotoBookingID uint64 `json:"photoBookingId,omitempty"`
}

// ListingAgendaDetailResponse wraps agenda entries for a listing.
type ListingAgendaDetailResponse struct {
	Entries    []ScheduleEntryResponse `json:"entries"`
	Pagination PaginationResponse      `json:"pagination"`
}

// ScheduleBlockEntryRequest carries data to create a blocking entry.
type ScheduleBlockEntryRequest struct {
	ListingID int64  `json:"listingId" binding:"required"`
	EntryType string `json:"entryType" binding:"required"`
	StartsAt  string `json:"startsAt" binding:"required"`
	EndsAt    string `json:"endsAt" binding:"required"`
	Reason    string `json:"reason,omitempty"`
}

// ScheduleBlockEntryUpdateRequest carries data to update a blocking entry.
type ScheduleBlockEntryUpdateRequest struct {
	EntryID   uint64 `json:"entryId" binding:"required"`
	ListingID int64  `json:"listingId" binding:"required"`
	EntryType string `json:"entryType" binding:"required"`
	StartsAt  string `json:"startsAt" binding:"required"`
	EndsAt    string `json:"endsAt" binding:"required"`
	Reason    string `json:"reason,omitempty"`
}

// ScheduleDeleteEntryRequest carries identifiers required to delete an agenda entry.
type ScheduleDeleteEntryRequest struct {
	EntryID   uint64 `json:"entryId" binding:"required"`
	ListingID int64  `json:"listingId" binding:"required"`
}

// ScheduleAvailabilityQuery represents query parameters to fetch listing availability slots.
type ScheduleAvailabilityQuery struct {
	ListingID          int64  `form:"listingId" binding:"required"`
	RangeFrom          string `form:"rangeFrom"`
	RangeTo            string `form:"rangeTo"`
	SlotDurationMinute uint16 `form:"slotDurationMinute"`
	Page               int    `form:"page"`
	Limit              int    `form:"limit"`
}

// ScheduleAvailabilitySlotResponse represents a continuous free window.
type ScheduleAvailabilitySlotResponse struct {
	StartsAt string `json:"startsAt"`
	EndsAt   string `json:"endsAt"`
}

// ScheduleAvailabilityResponse aggregates availability slots with pagination metadata.
type ScheduleAvailabilityResponse struct {
	Slots      []ScheduleAvailabilitySlotResponse `json:"slots"`
	Pagination PaginationResponse                 `json:"pagination"`
}

// ScheduleBlockEntryResponse represents a blocking entry returned by create/update operations.
type ScheduleBlockEntryResponse struct {
	Entry ScheduleEntryResponse `json:"entry"`
}
