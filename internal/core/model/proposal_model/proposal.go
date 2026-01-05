package proposalmodel

import (
	"database/sql"
	"time"
)

// ProposalInterface defines the contract for a proposal domain object.
type ProposalInterface interface {
	ID() int64
	SetID(id int64)
	ListingIdentityID() int64
	SetListingIdentityID(id int64)
	RealtorID() int64
	SetRealtorID(id int64)
	OwnerID() int64
	SetOwnerID(id int64)
	TransactionType() TransactionType
	SetTransactionType(t TransactionType)
	PaymentMethod() PaymentMethod
	SetPaymentMethod(m PaymentMethod)
	ProposedValue() float64
	SetProposedValue(v float64)
	OriginalValue() float64
	SetOriginalValue(v float64)
	DownPayment() sql.NullFloat64
	SetDownPayment(v sql.NullFloat64)
	Installments() sql.NullInt64
	SetInstallments(v sql.NullInt64)
	AcceptsExchange() bool
	SetAcceptsExchange(v bool)
	RentalMonths() sql.NullInt64
	SetRentalMonths(v sql.NullInt64)
	GuaranteeType() sql.NullString
	SetGuaranteeType(v sql.NullString)
	SecurityDeposit() sql.NullFloat64
	SetSecurityDeposit(v sql.NullFloat64)
	ClientName() string
	SetClientName(v string)
	ClientPhone() string
	SetClientPhone(v string)
	ProposalNotes() sql.NullString
	SetProposalNotes(v sql.NullString)
	OwnerNotes() sql.NullString
	SetOwnerNotes(v sql.NullString)
	RejectionReason() sql.NullString
	SetRejectionReason(v sql.NullString)
	Status() Status
	SetStatus(s Status)
	ExpiresAt() sql.NullTime
	SetExpiresAt(t sql.NullTime)
	AcceptedAt() sql.NullTime
	SetAcceptedAt(t sql.NullTime)
	RejectedAt() sql.NullTime
	SetRejectedAt(t sql.NullTime)
	CancelledAt() sql.NullTime
	SetCancelledAt(t sql.NullTime)
	CreatedAt() time.Time
	SetCreatedAt(t time.Time)
	UpdatedAt() time.Time
	SetUpdatedAt(t time.Time)
}

// Proposal is the concrete domain implementation of ProposalInterface.
type Proposal struct {
	id                int64
	listingIdentityID int64
	realtorID         int64
	ownerID           int64
	transactionType   TransactionType
	paymentMethod     PaymentMethod
	proposedValue     float64
	originalValue     float64
	downPayment       sql.NullFloat64
	installments      sql.NullInt64
	acceptsExchange   bool
	rentalMonths      sql.NullInt64
	guaranteeType     sql.NullString
	securityDeposit   sql.NullFloat64
	clientName        string
	clientPhone       string
	proposalNotes     sql.NullString
	ownerNotes        sql.NullString
	rejectionReason   sql.NullString
	status            Status
	expiresAt         sql.NullTime
	acceptedAt        sql.NullTime
	rejectedAt        sql.NullTime
	cancelledAt       sql.NullTime
	createdAt         time.Time
	updatedAt         time.Time
}

// NewProposal builds an empty proposal domain object.
func NewProposal() ProposalInterface {
	return &Proposal{}
}

func (p *Proposal) ID() int64      { return p.id }
func (p *Proposal) SetID(id int64) { p.id = id }

func (p *Proposal) ListingIdentityID() int64      { return p.listingIdentityID }
func (p *Proposal) SetListingIdentityID(id int64) { p.listingIdentityID = id }

func (p *Proposal) RealtorID() int64      { return p.realtorID }
func (p *Proposal) SetRealtorID(id int64) { p.realtorID = id }

func (p *Proposal) OwnerID() int64      { return p.ownerID }
func (p *Proposal) SetOwnerID(id int64) { p.ownerID = id }

func (p *Proposal) TransactionType() TransactionType     { return p.transactionType }
func (p *Proposal) SetTransactionType(t TransactionType) { p.transactionType = t }

func (p *Proposal) PaymentMethod() PaymentMethod     { return p.paymentMethod }
func (p *Proposal) SetPaymentMethod(m PaymentMethod) { p.paymentMethod = m }

func (p *Proposal) ProposedValue() float64     { return p.proposedValue }
func (p *Proposal) SetProposedValue(v float64) { p.proposedValue = v }

func (p *Proposal) OriginalValue() float64     { return p.originalValue }
func (p *Proposal) SetOriginalValue(v float64) { p.originalValue = v }

func (p *Proposal) DownPayment() sql.NullFloat64     { return p.downPayment }
func (p *Proposal) SetDownPayment(v sql.NullFloat64) { p.downPayment = v }

func (p *Proposal) Installments() sql.NullInt64     { return p.installments }
func (p *Proposal) SetInstallments(v sql.NullInt64) { p.installments = v }

func (p *Proposal) AcceptsExchange() bool     { return p.acceptsExchange }
func (p *Proposal) SetAcceptsExchange(v bool) { p.acceptsExchange = v }

func (p *Proposal) RentalMonths() sql.NullInt64     { return p.rentalMonths }
func (p *Proposal) SetRentalMonths(v sql.NullInt64) { p.rentalMonths = v }

func (p *Proposal) GuaranteeType() sql.NullString     { return p.guaranteeType }
func (p *Proposal) SetGuaranteeType(v sql.NullString) { p.guaranteeType = v }

func (p *Proposal) SecurityDeposit() sql.NullFloat64     { return p.securityDeposit }
func (p *Proposal) SetSecurityDeposit(v sql.NullFloat64) { p.securityDeposit = v }

func (p *Proposal) ClientName() string     { return p.clientName }
func (p *Proposal) SetClientName(v string) { p.clientName = v }

func (p *Proposal) ClientPhone() string     { return p.clientPhone }
func (p *Proposal) SetClientPhone(v string) { p.clientPhone = v }

func (p *Proposal) ProposalNotes() sql.NullString     { return p.proposalNotes }
func (p *Proposal) SetProposalNotes(v sql.NullString) { p.proposalNotes = v }

func (p *Proposal) OwnerNotes() sql.NullString     { return p.ownerNotes }
func (p *Proposal) SetOwnerNotes(v sql.NullString) { p.ownerNotes = v }

func (p *Proposal) RejectionReason() sql.NullString     { return p.rejectionReason }
func (p *Proposal) SetRejectionReason(v sql.NullString) { p.rejectionReason = v }

func (p *Proposal) Status() Status     { return p.status }
func (p *Proposal) SetStatus(s Status) { p.status = s }

func (p *Proposal) ExpiresAt() sql.NullTime     { return p.expiresAt }
func (p *Proposal) SetExpiresAt(t sql.NullTime) { p.expiresAt = t }

func (p *Proposal) AcceptedAt() sql.NullTime     { return p.acceptedAt }
func (p *Proposal) SetAcceptedAt(t sql.NullTime) { p.acceptedAt = t }

func (p *Proposal) RejectedAt() sql.NullTime     { return p.rejectedAt }
func (p *Proposal) SetRejectedAt(t sql.NullTime) { p.rejectedAt = t }

func (p *Proposal) CancelledAt() sql.NullTime     { return p.cancelledAt }
func (p *Proposal) SetCancelledAt(t sql.NullTime) { p.cancelledAt = t }

func (p *Proposal) CreatedAt() time.Time     { return p.createdAt }
func (p *Proposal) SetCreatedAt(t time.Time) { p.createdAt = t }

func (p *Proposal) UpdatedAt() time.Time     { return p.updatedAt }
func (p *Proposal) SetUpdatedAt(t time.Time) { p.updatedAt = t }
