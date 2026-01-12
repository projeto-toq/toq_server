package proposalmodel

import "database/sql"

// RealtorSummary represents aggregated metadata about the realtor that authored a proposal, including account age, accepted proposal counts and signed photo URLs.
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
	AcceptedProposals() int64
	SetAcceptedProposals(int64)
	PhotoURL() string
	SetPhotoURL(string)
}

type realtorSummary struct {
	id                int64
	name              string
	nickname          sql.NullString
	usageMonths       int
	proposalsCreated  int64
	acceptedProposals int64
	photoURL          string
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

func (r *realtorSummary) AcceptedProposals() int64 { return r.acceptedProposals }

func (r *realtorSummary) SetAcceptedProposals(total int64) { r.acceptedProposals = total }

func (r *realtorSummary) PhotoURL() string { return r.photoURL }

func (r *realtorSummary) SetPhotoURL(url string) { r.photoURL = url }
