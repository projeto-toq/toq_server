package entities

import "database/sql"

// RealtorSummaryEntity represents the aggregated realtor information fetched from users/proposals tables.
type RealtorSummaryEntity struct {
	RealtorID      int64
	FullName       string
	NickName       sql.NullString
	UsageMonths    sql.NullInt64
	ProposalsCount sql.NullInt64
}
