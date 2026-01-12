package entities

import "database/sql"

// OwnerMetricsEntity mirrors owner_response_metrics table for persistence operations.
type OwnerMetricsEntity struct {
	UserID int64

	VisitAvgResponseSeconds    sql.NullInt64
	VisitTotalResponses        int64
	VisitLastResponseAt        sql.NullTime
	ProposalAvgResponseSeconds sql.NullInt64
	ProposalTotalResponses     int64
	ProposalLastResponseAt     sql.NullTime
}
