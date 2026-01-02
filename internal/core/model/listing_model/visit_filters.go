package listingmodel

import "time"

// VisitListFilter constrains visit lookups for owners or requesters.
type VisitListFilter struct {
	ListingIdentityID *int64
	OwnerUserID       *int64
	RequesterUserID   *int64
	Statuses          []VisitStatus
	From              *time.Time
	To                *time.Time
	Page              int
	Limit             int
}

// VisitListResult holds a paginated visit collection.
type VisitListResult struct {
	Visits []VisitInterface
	Total  int64
}
