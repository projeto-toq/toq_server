package listingmodel

import (
	"database/sql"
	"strings"
	"time"
)

// VisitParticipantSnapshot stores immutable metadata about either the owner or the realtor tied to a visit.
//
// This snapshot is generated in the repository layer so the service can enrich API responses without performing
// additional queries per participant. It intentionally carries raw values (sql.NullInt64) so callers can decide how
// to present optional metrics like average response time or total visits performed.
type VisitParticipantSnapshot struct {
	// UserID identifies the participant (owner or realtor) associated with the visit.
	UserID int64

	// FullName mirrors users.full_name and is already trimmed in the repository.
	FullName string

	// CreatedAt indicates when the user account was created (users.created_at).
	CreatedAt time.Time

	// PhotoURL is populated by the service layer via userservices.GetPhotoDownloadURL for UI consumption.
	PhotoURL string

	// AvgResponseSeconds holds owner_response_metrics.visit_avg_response_time_seconds when available.
	AvgResponseSeconds sql.NullInt64

	// TotalVisits stores the total visits performed by the participant when applicable (realtor side only).
	TotalVisits sql.NullInt64
}

// SanitizedName trims excessive whitespace to keep downstream DTO builders lean.
func (s VisitParticipantSnapshot) SanitizedName() string {
	return strings.TrimSpace(s.FullName)
}

// HasAvgResponseSeconds reports whether AvgResponseSeconds carries a meaningful value.
func (s VisitParticipantSnapshot) HasAvgResponseSeconds() bool {
	return s.AvgResponseSeconds.Valid && s.AvgResponseSeconds.Int64 > 0
}

// AvgResponseSecondsValue returns the stored seconds or zero when absent. Callers should check HasAvgResponseSeconds.
func (s VisitParticipantSnapshot) AvgResponseSecondsValue() int64 {
	if !s.AvgResponseSeconds.Valid {
		return 0
	}
	return s.AvgResponseSeconds.Int64
}

// TotalVisitsValue returns the stored visit count or zero when not populated.
func (s VisitParticipantSnapshot) TotalVisitsValue() int64 {
	if !s.TotalVisits.Valid {
		return 0
	}
	return s.TotalVisits.Int64
}
