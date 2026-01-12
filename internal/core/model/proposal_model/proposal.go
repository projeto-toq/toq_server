package proposalmodel

import (
	"database/sql"
	"time"
)

// ProposalInterface models the attributes used by the new proposal flows.
type ProposalInterface interface {
	ID() int64
	SetID(int64)
	ListingIdentityID() int64
	SetListingIdentityID(int64)
	RealtorID() int64
	SetRealtorID(int64)
	OwnerID() int64
	SetOwnerID(int64)
	ProposalText() string
	SetProposalText(string)
	RejectionReason() sql.NullString
	SetRejectionReason(sql.NullString)
	Status() Status
	SetStatus(Status)
	AcceptedAt() sql.NullTime
	SetAcceptedAt(sql.NullTime)
	RejectedAt() sql.NullTime
	SetRejectedAt(sql.NullTime)
	CancelledAt() sql.NullTime
	SetCancelledAt(sql.NullTime)
	DocumentsCount() int
	SetDocumentsCount(int)
	CreatedAt() time.Time
	SetCreatedAt(time.Time)
	FirstOwnerActionAt() sql.NullTime
	SetFirstOwnerActionAt(sql.NullTime)
}

type proposal struct {
	id                 int64
	listingIdentityID  int64
	realtorID          int64
	ownerID            int64
	proposalText       string
	rejectionReason    sql.NullString
	status             Status
	acceptedAt         sql.NullTime
	rejectedAt         sql.NullTime
	cancelledAt        sql.NullTime
	documentsCount     int
	createdAt          time.Time
	firstOwnerActionAt sql.NullTime
}

// NewProposal instantiates an empty proposal domain object.
func NewProposal() ProposalInterface {
	return &proposal{}
}

func (p *proposal) ID() int64                { return p.id }
func (p *proposal) SetID(id int64)           { p.id = id }
func (p *proposal) ListingIdentityID() int64 { return p.listingIdentityID }
func (p *proposal) SetListingIdentityID(id int64) {
	p.listingIdentityID = id
}
func (p *proposal) RealtorID() int64 { return p.realtorID }
func (p *proposal) SetRealtorID(id int64) {
	p.realtorID = id
}
func (p *proposal) OwnerID() int64 { return p.ownerID }
func (p *proposal) SetOwnerID(id int64) {
	p.ownerID = id
}
func (p *proposal) ProposalText() string { return p.proposalText }
func (p *proposal) SetProposalText(text string) {
	p.proposalText = text
}
func (p *proposal) RejectionReason() sql.NullString { return p.rejectionReason }
func (p *proposal) SetRejectionReason(reason sql.NullString) {
	p.rejectionReason = reason
}
func (p *proposal) Status() Status { return p.status }
func (p *proposal) SetStatus(status Status) {
	p.status = status
}
func (p *proposal) AcceptedAt() sql.NullTime { return p.acceptedAt }
func (p *proposal) SetAcceptedAt(ts sql.NullTime) {
	p.acceptedAt = ts
}
func (p *proposal) RejectedAt() sql.NullTime { return p.rejectedAt }
func (p *proposal) SetRejectedAt(ts sql.NullTime) {
	p.rejectedAt = ts
}
func (p *proposal) CancelledAt() sql.NullTime { return p.cancelledAt }
func (p *proposal) SetCancelledAt(ts sql.NullTime) {
	p.cancelledAt = ts
}
func (p *proposal) DocumentsCount() int { return p.documentsCount }
func (p *proposal) SetDocumentsCount(count int) {
	p.documentsCount = count
}
func (p *proposal) CreatedAt() time.Time { return p.createdAt }
func (p *proposal) SetCreatedAt(ts time.Time) {
	p.createdAt = ts
}
func (p *proposal) FirstOwnerActionAt() sql.NullTime { return p.firstOwnerActionAt }
func (p *proposal) SetFirstOwnerActionAt(ts sql.NullTime) {
	p.firstOwnerActionAt = ts
}
