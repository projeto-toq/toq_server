package proposalmodel

import (
	"database/sql"
	"time"
)

// Stats aggregates proposal counters and monetary extremes for reporting.
type Stats struct {
	TotalProposals  int64           // total proposals in the queried scope
	PendingCount    int64           // pending proposals
	AcceptedCount   int64           // accepted proposals
	RejectedCount   int64           // rejected proposals
	CancelledCount  int64           // cancelled proposals
	ExpiredCount    int64           // expired proposals
	HighestProposal sql.NullFloat64 // highest proposed value
	LowestProposal  sql.NullFloat64 // lowest proposed value
	AverageProposal sql.NullFloat64 // average proposed value
}

// StatsFilter narrows stats aggregation for specific users or time ranges.
type StatsFilter struct {
	RealtorID *int64     // optional realtor scope
	OwnerID   *int64     // optional owner scope
	StartDate *time.Time // optional start datetime filter
	EndDate   *time.Time // optional end datetime filter
}
