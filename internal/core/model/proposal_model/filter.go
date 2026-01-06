package proposalmodel

// ActorScope narrows list queries based on the authenticated role.
type ActorScope string

const (
	ActorScopeRealtor ActorScope = "realtor"
	ActorScopeOwner   ActorScope = "owner"
)

// ListFilter stores normalized filters for repository queries.
type ListFilter struct {
	ActorScope ActorScope
	ActorID    int64
	ListingID  *int64
	Statuses   []Status
	Page       int
	Limit      int
}

// ListResult bundles the paginated proposals and the total counter.
type ListResult struct {
	Items []ProposalInterface
	Total int64
}
