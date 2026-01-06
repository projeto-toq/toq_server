package entities

import "database/sql"

// ProposalEntity mirrors the proposals table schema for persistence operations.
type ProposalEntity struct {
	ID                int64
	ListingIdentityID int64
	RealtorID         int64
	OwnerID           int64
	ProposalText      sql.NullString
	RejectionReason   sql.NullString
	Status            string
	AcceptedAt        sql.NullTime
	RejectedAt        sql.NullTime
	CancelledAt       sql.NullTime
	Deleted           bool
	DocumentsCount    sql.NullInt64
}
