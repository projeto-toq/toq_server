package entities

import (
	"database/sql"
	"time"
)

// ProposalEntity mirrors the proposals table schema for persistence operations.
type ProposalEntity struct {
	ID                int64
	ListingIdentityID int64
	RealtorID         int64
	OwnerID           int64
	TransactionType   string
	PaymentMethod     string
	ProposedValue     float64
	OriginalValue     float64
	DownPayment       sql.NullFloat64
	Installments      sql.NullInt64
	AcceptsExchange   bool
	RentalMonths      sql.NullInt64
	GuaranteeType     sql.NullString
	SecurityDeposit   sql.NullFloat64
	ClientName        string
	ClientPhone       string
	ProposalNotes     sql.NullString
	OwnerNotes        sql.NullString
	RejectionReason   sql.NullString
	Status            string
	ExpiresAt         sql.NullTime
	AcceptedAt        sql.NullTime
	RejectedAt        sql.NullTime
	CancelledAt       sql.NullTime
	IsFavorite        bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Deleted           bool
}
