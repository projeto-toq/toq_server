package proposalmodel

import "database/sql"

// OwnerSummary represents aggregated metadata about the listing owner for proposal views.
// It exposes basic identity, engagement metrics and media fields while keeping the domain
// decoupled from HTTP/DB concerns.
type OwnerSummary interface {
	ID() int64
	SetID(int64)

	FullName() string
	SetFullName(string)

	MemberSinceMonths() int
	SetMemberSinceMonths(int)

	PhotoURL() string
	SetPhotoURL(string)

	ProposalAvgSeconds() sql.NullInt64
	SetProposalAvgSeconds(sql.NullInt64)

	VisitAvgSeconds() sql.NullInt64
	SetVisitAvgSeconds(sql.NullInt64)
}

type ownerSummary struct {
	id                int64
	fullName          string
	memberSinceMonths int
	photoURL          string
	proposalAvg       sql.NullInt64
	visitAvg          sql.NullInt64
}

// NewOwnerSummary instantiates an empty owner summary domain object.
func NewOwnerSummary() OwnerSummary {
	return &ownerSummary{}
}

func (o *ownerSummary) ID() int64                         { return o.id }
func (o *ownerSummary) SetID(id int64)                    { o.id = id }
func (o *ownerSummary) FullName() string                  { return o.fullName }
func (o *ownerSummary) SetFullName(name string)           { o.fullName = name }
func (o *ownerSummary) MemberSinceMonths() int            { return o.memberSinceMonths }
func (o *ownerSummary) SetMemberSinceMonths(months int)   { o.memberSinceMonths = months }
func (o *ownerSummary) PhotoURL() string                  { return o.photoURL }
func (o *ownerSummary) SetPhotoURL(url string)            { o.photoURL = url }
func (o *ownerSummary) ProposalAvgSeconds() sql.NullInt64 { return o.proposalAvg }
func (o *ownerSummary) SetProposalAvgSeconds(v sql.NullInt64) {
	o.proposalAvg = v
}
func (o *ownerSummary) VisitAvgSeconds() sql.NullInt64 { return o.visitAvg }
func (o *ownerSummary) SetVisitAvgSeconds(v sql.NullInt64) {
	o.visitAvg = v
}
