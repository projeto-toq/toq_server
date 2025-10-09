package devicetokenrepository

import usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"

// DeviceTokenRepoPortInterface defines persistence operations for device tokens (hexagonal right port)
type DeviceTokenRepoPortInterface interface {
	ListByUserID(userID int64) ([]usermodel.DeviceTokenInterface, error)
	AddToken(userID int64, token string, platform *string) (usermodel.DeviceTokenInterface, error)
	RemoveToken(userID int64, token string) error
	RemoveAllByUserID(userID int64) error
	ListTokensByOptedInUsers() ([]string, error)
	ListTokensByUserIDIfOptedIn(userID int64) ([]string, error)

	// Per-device operations (compatible defaults when device_id is not present in schema)
	AddTokenForDevice(userID int64, deviceID, token string, platform *string) (usermodel.DeviceTokenInterface, error)
	RemoveTokensByDeviceID(userID int64, deviceID string) error
	ListTokensByDeviceID(userID int64, deviceID string) ([]string, error)
}
