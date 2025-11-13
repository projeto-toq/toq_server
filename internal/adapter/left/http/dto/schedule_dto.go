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

// ScheduleRuleRequest represents the payload to create recurring unavailability rules.
type ScheduleRuleRequest struct {
	ListingIdentityID int64    `json:"listingIdentityId" binding:"required" example:"3241"`
	WeekDays          []string `json:"weekDays" binding:"required" example:"[\"MONDAY\",\"TUESDAY\"]"`
	RangeStart        string   `json:"rangeStart" binding:"required" example:"09:00"`
	RangeEnd          string   `json:"rangeEnd" binding:"required" example:"18:00"`
	Active            bool     `json:"active"`
	Timezone          string   `json:"timezone" binding:"required" example:"America/Sao_Paulo"`
}

// ScheduleRuleUpdateRequest represents the payload to update an existing rule.
type ScheduleRuleUpdateRequest struct {
	RuleID            uint64   `json:"ruleId" binding:"required" example:"9801"`
	ListingIdentityID int64    `json:"listingIdentityId" binding:"required" example:"3241"`
	WeekDays          []string `json:"weekDays" binding:"required" example:"[\"MONDAY\"]"`
	RangeStart        string   `json:"rangeStart" binding:"required" example:"10:00"`
	RangeEnd          string   `json:"rangeEnd" binding:"required" example:"22:00"`
	Active            bool     `json:"active"`
	Timezone          string   `json:"timezone" binding:"required" example:"America/Sao_Paulo"`
}

// ScheduleRuleDeleteRequest represents the payload to delete a recurring rule.
type ScheduleRuleDeleteRequest struct {
	RuleID            uint64 `json:"ruleId" binding:"required" example:"9801"`
	ListingIdentityID int64  `json:"listingIdentityId" binding:"required" example:"3241"`
}

// ScheduleRuleResponse exposes a recurring rule definition.
type ScheduleRuleResponse struct {
	RuleID    uint64 `json:"ruleId"`
	Weekday   string `json:"weekday"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Active    bool   `json:"active"`
}

// ScheduleRulesResponse wraps a listing rule collection.
type ScheduleRulesResponse struct {
	ListingIdentityID int64                  `json:"listingIdentityId"`
	Rules             []ScheduleRuleResponse `json:"rules"`
	Timezone          string                 `json:"timezone"`
}

// ScheduleFinishAgendaRequest represents the payload to finish agenda creation.
type ScheduleFinishAgendaRequest struct {
	ListingIdentityID int64 `json:"listingIdentityId" binding:"required" example:"3241"`
}

// OwnerAgendaSummaryQuery captures query string parameters for owner agenda summary.
type OwnerAgendaSummaryQuery struct {
	ListingIdentityIDs []int64 `form:"listingIdentityIds"`
	RangeFrom          string  `form:"rangeFrom"`
	RangeTo            string  `form:"rangeTo"`
	Page               int     `form:"page"`
	Limit              int     `form:"limit"`
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
	ListingIdentityID int64                             `json:"listingIdentityId"`
	Entries           []OwnerAgendaSummaryEntryResponse `json:"entries"`
}

// OwnerAgendaSummaryResponse aggregates the consolidated agenda view for owners.
type OwnerAgendaSummaryResponse struct {
	Items      []OwnerAgendaSummaryItemResponse `json:"items"`
	Pagination PaginationResponse               `json:"pagination"`
}

// ListingAgendaDetailQuery represents query parameters to list agenda entries of a specific listing.
type ListingAgendaDetailQuery struct {
	ListingIdentityID int64  `form:"listingIdentityId" binding:"required"`
	RangeFrom         string `form:"rangeFrom"`
	RangeTo           string `form:"rangeTo"`
	Page              int    `form:"page"`
	Limit             int    `form:"limit"`
}

// ScheduleEntryResponse exposes detailed information about a single agenda entry.
type ScheduleEntryResponse struct {
	ID             uint64 `json:"id,omitempty"`
	RuleID         uint64 `json:"ruleId,omitempty"`
	SourceType     string `json:"sourceType"`
	Recurring      bool   `json:"recurring"`
	Weekday        string `json:"weekday,omitempty"`
	EntryType      string `json:"entryType,omitempty"`
	StartsAt       string `json:"startsAt"`
	EndsAt         string `json:"endsAt"`
	Blocking       bool   `json:"blocking"`
	Reason         string `json:"reason,omitempty"`
	VisitID        uint64 `json:"visitId,omitempty"`
	PhotoBookingID uint64 `json:"photoBookingId,omitempty"`
	Timezone       string `json:"timezone"`
}

// ListingAgendaDetailResponse wraps agenda entries for a listing.
type ListingAgendaDetailResponse struct {
	Entries    []ScheduleEntryResponse `json:"entries"`
	Pagination PaginationResponse      `json:"pagination"`
	Timezone   string                  `json:"timezone"`
}

// ListingBlockRulesQuery represents query parameters to list recurring block rules for a listing agenda.
type ListingBlockRulesQuery struct {
	ListingIdentityID int64    `form:"listingIdentityId" binding:"required"`
	WeekDays          []string `form:"weekDays"`
}

// ScheduleAvailabilityQuery represents query parameters to fetch listing availability slots.
type ScheduleAvailabilityQuery struct {
	ListingIdentityID  int64  `form:"listingIdentityId" binding:"required"`
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
	Timezone   string                             `json:"timezone"`
}
