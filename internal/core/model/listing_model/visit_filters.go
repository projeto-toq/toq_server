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

// VisitWithListing aggregates a visit and its active listing snapshot.
type VisitWithListing struct {
	Visit   VisitInterface
	Listing ListingInterface
}

// VisitListResult holds a paginated visit collection with listing snapshots attached.
type VisitListResult struct {
	Visits []VisitWithListing
	Total  int64
}
