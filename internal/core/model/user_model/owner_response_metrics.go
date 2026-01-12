package usermodel

import "database/sql"

// OwnerResponseMetrics models aggregated SLA indicators for owner interactions.
type OwnerResponseMetrics interface {
	OwnerID() int64
	SetOwnerID(int64)

	VisitAverageSeconds() sql.NullInt64
	SetVisitAverageSeconds(sql.NullInt64)
	VisitResponsesTotal() int64
	SetVisitResponsesTotal(int64)
	VisitLastResponseAt() sql.NullTime
	SetVisitLastResponseAt(sql.NullTime)

	ProposalAverageSeconds() sql.NullInt64
	SetProposalAverageSeconds(sql.NullInt64)
	ProposalResponsesTotal() int64
	SetProposalResponsesTotal(int64)
	ProposalLastResponseAt() sql.NullTime
	SetProposalLastResponseAt(sql.NullTime)
}

type ownerResponseMetrics struct {
	ownerID int64

	visitAvgSeconds sql.NullInt64
	visitResponses  int64
	visitLastAt     sql.NullTime

	proposalAvgSeconds sql.NullInt64
	proposalResponses  int64
	proposalLastAt     sql.NullTime
}

// NewOwnerResponseMetrics instantiates an empty metrics aggregate.
func NewOwnerResponseMetrics() OwnerResponseMetrics {
	return &ownerResponseMetrics{}
}

func (m *ownerResponseMetrics) OwnerID() int64                     { return m.ownerID }
func (m *ownerResponseMetrics) SetOwnerID(id int64)                { m.ownerID = id }
func (m *ownerResponseMetrics) VisitAverageSeconds() sql.NullInt64 { return m.visitAvgSeconds }
func (m *ownerResponseMetrics) SetVisitAverageSeconds(value sql.NullInt64) {
	m.visitAvgSeconds = value
}
func (m *ownerResponseMetrics) VisitResponsesTotal() int64 { return m.visitResponses }
func (m *ownerResponseMetrics) SetVisitResponsesTotal(total int64) {
	m.visitResponses = total
}
func (m *ownerResponseMetrics) VisitLastResponseAt() sql.NullTime { return m.visitLastAt }
func (m *ownerResponseMetrics) SetVisitLastResponseAt(ts sql.NullTime) {
	m.visitLastAt = ts
}
func (m *ownerResponseMetrics) ProposalAverageSeconds() sql.NullInt64 { return m.proposalAvgSeconds }
func (m *ownerResponseMetrics) SetProposalAverageSeconds(value sql.NullInt64) {
	m.proposalAvgSeconds = value
}
func (m *ownerResponseMetrics) ProposalResponsesTotal() int64 { return m.proposalResponses }
func (m *ownerResponseMetrics) SetProposalResponsesTotal(total int64) {
	m.proposalResponses = total
}
func (m *ownerResponseMetrics) ProposalLastResponseAt() sql.NullTime { return m.proposalLastAt }
func (m *ownerResponseMetrics) SetProposalLastResponseAt(ts sql.NullTime) {
	m.proposalLastAt = ts
}
