package proposalmodel

import "database/sql"

// RealtorSummary represents aggregated metadata about the realtor that authored a proposal.
type RealtorSummary interface {
	ID() int64
	SetID(int64)
	Name() string
	SetName(string)
	Nickname() sql.NullString
	SetNickname(sql.NullString)
	UsageMonths() int
	SetUsageMonths(int)
	ProposalsCreated() int64
	SetProposalsCreated(int64)
}

type realtorSummary struct {
	id               int64
	name             string
	nickname         sql.NullString
	usageMonths      int
	proposalsCreated int64
}

// NewRealtorSummary instantiates an empty realtor summary object.
func NewRealtorSummary() RealtorSummary {
	return &realtorSummary{}
}

func (r *realtorSummary) ID() int64 { return r.id }

func (r *realtorSummary) SetID(id int64) { r.id = id }

func (r *realtorSummary) Name() string { return r.name }

func (r *realtorSummary) SetName(name string) { r.name = name }

func (r *realtorSummary) Nickname() sql.NullString { return r.nickname }

func (r *realtorSummary) SetNickname(nick sql.NullString) { r.nickname = nick }

func (r *realtorSummary) UsageMonths() int { return r.usageMonths }

func (r *realtorSummary) SetUsageMonths(months int) { r.usageMonths = months }

func (r *realtorSummary) ProposalsCreated() int64 { return r.proposalsCreated }

func (r *realtorSummary) SetProposalsCreated(total int64) { r.proposalsCreated = total }
