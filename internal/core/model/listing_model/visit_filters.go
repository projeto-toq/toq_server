package listingmodel

import "time"

// VisitListFilter constrains visit lookups for owners or realtors.
type VisitListFilter struct {
	ListingID *int64
	OwnerID   *int64
	RealtorID *int64
	Statuses  []VisitStatus
	From      *time.Time
	To        *time.Time
	Page      int
	Limit     int
}

// VisitListResult holds a paginated visit collection.
type VisitListResult struct {
	Visits []VisitInterface
	Total  int64
}
