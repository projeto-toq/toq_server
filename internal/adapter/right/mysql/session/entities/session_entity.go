package sessionentities

import "time"

type SessionEntity struct {
	ID                int64
	UserID            int64
	RefreshHash       string
	TokenJTI          string
	ExpiresAt         time.Time
	AbsoluteExpiresAt time.Time
	CreatedAt         time.Time
	RotatedAt         *time.Time
	UserAgent         string
	IP                string
	DeviceID          string
	RotationCounter   int
	LastRefreshAt     *time.Time
	Revoked           bool
}
