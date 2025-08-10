package usermodel

// DeviceTokenInterface defines the access methods for a device push token record
// aligning with the existing pattern used for other domain interfaces in user_model.
type DeviceTokenInterface interface {
	GetID() int64
	SetID(int64)
	GetUserID() int64
	SetUserID(int64)
	GetDeviceToken() string
	SetDeviceToken(string)
}

// NewDeviceToken creates a new deviceToken implementing DeviceTokenInterface.
func NewDeviceToken() DeviceTokenInterface {
	return &deviceToken{}
}
