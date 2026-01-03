package sessionmodel

import "time"

// SessionInterface defines the contract for authentication sessions in the domain layer.
// Responsibilities:
//   - Represent refresh-session state without coupling to DB/HTTP concerns.
//   - Provide getter/setter API for services and repositories to read/update fields.
// Usage: services orchestrate lifecycle; repositories persist via converters; handlers should not import this type directly.
type SessionInterface interface {
	// GetID returns the session primary key (0 when not persisted).
	GetID() int64
	// SetID sets the session primary key (called by repository after insert).
	SetID(int64)

	// GetUserID returns the owner user id.
	GetUserID() int64
	// SetUserID sets the owner user id.
	SetUserID(int64)

	// GetRefreshHash returns the refresh token hash (SHA-256 hex, len=64).
	GetRefreshHash() string
	// SetRefreshHash sets the refresh token hash.
	SetRefreshHash(string)

	// GetTokenJTI returns the access token JTI (UUID) or empty when not issued.
	GetTokenJTI() string
	// SetTokenJTI sets the access token JTI.
	SetTokenJTI(string)

	// GetExpiresAt returns the sliding expiration timestamp.
	GetExpiresAt() time.Time
	// SetExpiresAt sets the sliding expiration timestamp.
	SetExpiresAt(time.Time)

	// GetAbsoluteExpiresAt returns the hard expiration timestamp (zero if none).
	GetAbsoluteExpiresAt() time.Time
	// SetAbsoluteExpiresAt sets the hard expiration timestamp.
	SetAbsoluteExpiresAt(time.Time)

	// GetCreatedAt returns when the session was created.
	GetCreatedAt() time.Time
	// SetCreatedAt sets the creation time.
	SetCreatedAt(time.Time)

	// GetRotatedAt returns last rotation timestamp or nil.
	GetRotatedAt() *time.Time
	// SetRotatedAt sets last rotation timestamp (nil to clear).
	SetRotatedAt(*time.Time)

	// GetUserAgent returns stored user agent (empty if none).
	GetUserAgent() string
	// SetUserAgent sets user agent string.
	SetUserAgent(string)

	// GetIP returns stored client IP (empty if none).
	GetIP() string
	// SetIP sets client IP.
	SetIP(string)

	// GetDeviceID returns device identifier (empty if none).
	GetDeviceID() string
	// SetDeviceID sets device identifier.
	SetDeviceID(string)

	// GetRotationCounter returns refresh rotation count.
	GetRotationCounter() int
	// SetRotationCounter sets refresh rotation count.
	SetRotationCounter(int)

	// GetLastRefreshAt returns last refresh timestamp or nil.
	GetLastRefreshAt() *time.Time
	// SetLastRefreshAt sets last refresh timestamp (nil to clear).
	SetLastRefreshAt(*time.Time)

	// GetRevoked indicates if session is revoked.
	GetRevoked() bool
	// SetRevoked sets revoked flag.
	SetRevoked(bool)
}

// Session is the concrete, unexported implementation of SessionInterface.
// It holds simple fields with getter/setter semantics to enforce encapsulation.
type Session struct {
	id          int64
	userID      int64
	refreshHash string
	tokenJTI    string
	expiresAt   time.Time
	absoluteExp time.Time
	createdAt   time.Time
	rotatedAt   *time.Time
	userAgent   string
	ip          string
	deviceID    string
	rotationCtr int
	lastRefresh *time.Time
	revoked     bool
}

// NewSession creates a new session with CreatedAt initialized to UTC now and default zero values elsewhere.
func NewSession() *Session {
	return &Session{createdAt: time.Now().UTC()}
}

func (s *Session) GetID() int64                     { return s.id }
func (s *Session) SetID(v int64)                    { s.id = v }
func (s *Session) GetUserID() int64                 { return s.userID }
func (s *Session) SetUserID(v int64)                { s.userID = v }
func (s *Session) GetRefreshHash() string           { return s.refreshHash }
func (s *Session) SetRefreshHash(v string)          { s.refreshHash = v }
func (s *Session) GetTokenJTI() string              { return s.tokenJTI }
func (s *Session) SetTokenJTI(v string)             { s.tokenJTI = v }
func (s *Session) GetExpiresAt() time.Time          { return s.expiresAt }
func (s *Session) SetExpiresAt(v time.Time)         { s.expiresAt = v }
func (s *Session) GetAbsoluteExpiresAt() time.Time  { return s.absoluteExp }
func (s *Session) SetAbsoluteExpiresAt(v time.Time) { s.absoluteExp = v }
func (s *Session) GetCreatedAt() time.Time          { return s.createdAt }
func (s *Session) SetCreatedAt(v time.Time)         { s.createdAt = v }
func (s *Session) GetRotatedAt() *time.Time         { return s.rotatedAt }
func (s *Session) SetRotatedAt(v *time.Time)        { s.rotatedAt = v }
func (s *Session) GetUserAgent() string             { return s.userAgent }
func (s *Session) SetUserAgent(v string)            { s.userAgent = v }
func (s *Session) GetIP() string                    { return s.ip }
func (s *Session) SetIP(v string)                   { s.ip = v }
func (s *Session) GetDeviceID() string              { return s.deviceID }
func (s *Session) SetDeviceID(v string)             { s.deviceID = v }
func (s *Session) GetRotationCounter() int          { return s.rotationCtr }
func (s *Session) SetRotationCounter(v int)         { s.rotationCtr = v }
func (s *Session) GetLastRefreshAt() *time.Time     { return s.lastRefresh }
func (s *Session) SetLastRefreshAt(v *time.Time)    { s.lastRefresh = v }
func (s *Session) GetRevoked() bool                 { return s.revoked }
func (s *Session) SetRevoked(v bool)                { s.revoked = v }
