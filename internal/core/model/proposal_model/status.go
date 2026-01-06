package proposalmodel

// Status represents the lifecycle state of a proposal.
type Status string

const (
	StatusPending   Status = "pending"
	StatusAccepted  Status = "accepted"
	StatusRefused   Status = "refused"
	StatusCancelled Status = "cancelled"
)

// String returns the textual representation of the status.
func (s Status) String() string { return string(s) }
