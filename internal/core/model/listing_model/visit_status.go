package listingmodel

import (
	"fmt"
	"strings"
)

// VisitStatus describes the lifecycle state of a property visit.
type VisitStatus string

const (
	VisitStatusPending   VisitStatus = "PENDING"
	VisitStatusApproved  VisitStatus = "APPROVED"
	VisitStatusRejected  VisitStatus = "REJECTED"
	VisitStatusCancelled VisitStatus = "CANCELLED"
	VisitStatusCompleted VisitStatus = "COMPLETED"
	VisitStatusNoShow    VisitStatus = "NO_SHOW"
)

// ParseVisitStatus converts a string to VisitStatus with case-insensitive matching.
func ParseVisitStatus(raw string) (VisitStatus, error) {
	normalized := strings.ToUpper(strings.TrimSpace(raw))
	switch normalized {
	case "PENDING":
		return VisitStatusPending, nil
	case "APPROVED":
		return VisitStatusApproved, nil
	case "REJECTED":
		return VisitStatusRejected, nil
	case "CANCELLED":
		return VisitStatusCancelled, nil
	case "COMPLETED":
		return VisitStatusCompleted, nil
	case "NO_SHOW":
		return VisitStatusNoShow, nil
	default:
		return "", fmt.Errorf("invalid visit status: %s", raw)
	}
}

// IsBlocking returns true when a visit status should block scheduling windows.
func (s VisitStatus) IsBlocking() bool {
	return s == VisitStatusPending || s == VisitStatusApproved
}
