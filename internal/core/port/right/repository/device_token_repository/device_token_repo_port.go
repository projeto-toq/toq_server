package devicetokenrepository

import usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

// DeviceTokenRepoPortInterface defines persistence operations for device tokens (hexagonal right port)
type DeviceTokenRepoPortInterface interface {
	ListByUserID(userID int64) ([]usermodel.DeviceTokenInterface, error)
	AddToken(userID int64, token string, platform *string) (usermodel.DeviceTokenInterface, error)
	RemoveToken(userID int64, token string) error
	RemoveAllByUserID(userID int64) error
}
