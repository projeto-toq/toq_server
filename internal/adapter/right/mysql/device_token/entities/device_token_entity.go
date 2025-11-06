package entities

import "database/sql"

// DeviceTokenEntity represents a row in the device_tokens table
type DeviceTokenEntity struct {
	ID       int64
	UserID   int64
	Token    string
	DeviceID string
	Platform sql.NullString
}
