package proposalmodel

import "time"

// ListFilter represents query parameters for listing proposals with pagination and ordering.
type ListFilter struct {
	Page              int               // 1-based page number
	Limit             int               // items per page (max enforced in handler/service)
	SortBy            string            // allowed fields: createdAt, proposedValue, expiresAt, status
	SortOrder         string            // asc or desc
	Statuses          []Status          // optional status filter
	ListingIdentityID *int64            // optional listing identity filter
	RealtorID         *int64            // optional realtor filter
	OwnerID           *int64            // optional owner filter
	StartDate         *time.Time        // inclusive start date filter (created_at)
	EndDate           *time.Time        // inclusive end date filter (created_at)
	MinValue          *float64          // minimum proposed value
	MaxValue          *float64          // maximum proposed value
	TransactionTypes  []TransactionType // optional transaction type filter
	PaymentMethods    []PaymentMethod   // optional payment method filter
	IncludeDeleted    bool              // include soft-deleted proposals when true
}

// ListResult is the unified output for list queries.
type ListResult struct {
	Items   []ProposalInterface // page items
	Total   int64               // total items for pagination
	Summary Stats               // aggregated summary (counts and monetary metrics)
}
