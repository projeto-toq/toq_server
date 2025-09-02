package events

// EventType enumerates supported event types in the system
type EventType string

const (
	SessionCreated  EventType = "session.created"
	SessionRotated  EventType = "session.rotated"
	SessionsRevoked EventType = "sessions.revoked"
)

// SessionEvent carries information about session lifecycle changes
type SessionEvent struct {
	Type      EventType
	UserID    int64
	SessionID *int64 // optional: present for create/rotate
	DeviceID  string // optional: for device-targeted actions
}
