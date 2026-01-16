package audientities

import "time"

// AuditEventEntity maps to the audit_events table columns.
type AuditEventEntity struct {
	ID             int64
	OccurredAt     time.Time
	ActorID        int64
	ActorRole      string
	ActorDeviceID  string
	ActorIP        string
	ActorUserAgent string
	TargetType     string
	TargetID       int64
	TargetVersion  int64
	Operation      string
	Metadata       any
	RequestID      string
	TraceID        string
}
