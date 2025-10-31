package schedulemodel

import "time"

// ScheduleRange defines the time window used to query agenda data.
type ScheduleRange struct {
	From time.Time
	To   time.Time
	Loc  *time.Location
}

// PaginationConfig groups pagination parameters shared by queries.
type PaginationConfig struct {
	Page  int
	Limit int
}

// OwnerSummaryFilter constrains consolidated agenda lookups for owners.
type OwnerSummaryFilter struct {
	OwnerID    int64
	ListingIDs []int64
	Range      ScheduleRange
	Pagination PaginationConfig
}

// AgendaDetailFilter constrains detailed agenda lookups for a single listing.
type AgendaDetailFilter struct {
	OwnerID    int64
	ListingID  int64
	Range      ScheduleRange
	Pagination PaginationConfig
}

// BlockRulesFilter constrains block rule lookups for a single listing.
type BlockRulesFilter struct {
	OwnerID   int64
	ListingID int64
	Weekdays  []time.Weekday
}

// AvailabilityFilter describes the parameters required to list free slots.
type AvailabilityFilter struct {
	ListingID          int64
	Range              ScheduleRange
	SlotDurationMinute uint16
	Pagination         PaginationConfig
}

// SummaryEntry represents a normalized entry shape for consolidated outputs.
type SummaryEntry struct {
	EntryType EntryType
	StartsAt  time.Time
	EndsAt    time.Time
	Blocking  bool
}

// OwnerSummaryItem groups entries per listing for owner dashboards.
type OwnerSummaryItem struct {
	ListingID int64
	Entries   []SummaryEntry
}

// OwnerSummaryResult captures the output produced by summary queries.
type OwnerSummaryResult struct {
	Items []OwnerSummaryItem
	Total int64
}

// AgendaEntriesPage represents a paginated collection retrieved from storage.
type AgendaEntriesPage struct {
	Entries []AgendaEntryInterface
	Total   int64
}

// AgendaTimelineItem groups concrete timeline windows for agenda visualisation.
type AgendaTimelineItem struct {
	Source    TimelineSource
	Entry     AgendaEntryInterface
	Rule      AgendaRuleInterface
	StartsAt  time.Time
	EndsAt    time.Time
	Weekday   time.Weekday
	Recurring bool
	Blocking  bool
}

// AgendaDetailResult contains paginated detail results already enriched for the handler.
type AgendaDetailResult struct {
	Items    []AgendaTimelineItem
	Total    int64
	Timezone string
}

// RuleListResult captures all rules configured for an agenda.
type RuleListResult struct {
	ListingID int64
	Timezone  string
	Rules     []AgendaRuleInterface
}

// AvailabilityData groups raw entries and rules so services can compute free slots.
type AvailabilityData struct {
	Entries []AgendaEntryInterface
	Rules   []AgendaRuleInterface
}
