package entities

import "database/sql"

// OwnerSummaryEntity represents aggregated owner data joined from users and owner_response_metrics.
type OwnerSummaryEntity struct {
	OwnerID            int64
	FullName           string
	MemberSinceMonths  sql.NullInt64
	ProposalAvgSeconds sql.NullInt64
	VisitAvgSeconds    sql.NullInt64
}
