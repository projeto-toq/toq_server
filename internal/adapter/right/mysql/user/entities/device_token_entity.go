package userentity

// DeviceTokenEntity represents a row in device_tokens table.
type DeviceTokenEntity struct {
	ID       int64
	UserID   int64
	Token    string
	DeviceID string
	Platform *string
}
