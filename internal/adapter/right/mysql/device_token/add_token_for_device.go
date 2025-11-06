package mysqldevicetokenadapter

import (
	devicetokenconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/device_token/converters"
	devicetokenentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/device_token/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// AddTokenForDevice adds or updates a token for a specific device
func (a *DeviceTokenAdapter) AddTokenForDevice(userID int64, deviceID, token string, platform *string) (usermodel.DeviceTokenInterface, error) {
	query := `INSERT INTO device_tokens (user_id, device_token, device_id, platform) 
  VALUES (?, ?, ?, ?) 
  ON DUPLICATE KEY UPDATE device_token = VALUES(device_token), platform = VALUES(platform)`

	var platformArg interface{}
	if platform != nil {
		platformArg = *platform
	}

	result, err := a.db.GetDB().Exec(query, userID, token, deviceID, platformArg)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	entity := devicetokenentity.DeviceTokenEntity{
		ID:       id,
		UserID:   userID,
		Token:    token,
		DeviceID: deviceID,
	}

	return devicetokenconverters.DeviceTokenEntityToDomain(entity), nil
}
