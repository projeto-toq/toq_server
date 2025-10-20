package listingmodel

// VisitStatus describes the lifecycle state of a property visit.
type VisitStatus string

const (
	VisitStatusPendingOwner VisitStatus = "PENDING_OWNER"
	VisitStatusConfirmed    VisitStatus = "CONFIRMED"
	VisitStatusCancelled    VisitStatus = "CANCELLED"
	VisitStatusDone         VisitStatus = "DONE"
)
