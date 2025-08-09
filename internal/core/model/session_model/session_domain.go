package sessionmodel

import "time"

type SessionInterface interface {
	GetID() int64
	SetID(int64)
	GetUserID() int64
	SetUserID(int64)
	GetRefreshHash() string
	SetRefreshHash(string)
	GetTokenJTI() string
	SetTokenJTI(string)
	GetExpiresAt() time.Time
	SetExpiresAt(time.Time)
	GetAbsoluteExpiresAt() time.Time
	SetAbsoluteExpiresAt(time.Time)
	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)
	GetRotatedAt() *time.Time
	SetRotatedAt(*time.Time)
	GetUserAgent() string
	SetUserAgent(string)
	GetIP() string
	SetIP(string)
	GetDeviceID() string
	SetDeviceID(string)
	GetRotationCounter() int
	SetRotationCounter(int)
	GetLastRefreshAt() *time.Time
	SetLastRefreshAt(*time.Time)
	GetRevoked() bool
	SetRevoked(bool)
}

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
