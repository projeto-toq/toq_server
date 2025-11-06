package mysqldevicetokenadapter

import (
	"database/sql"

	devicetokenconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/device_token/converters"
	devicetokenentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/device_token/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// AddToken adds a device token using INSERT IGNORE (legacy behavior)
func (a *DeviceTokenAdapter) AddToken(userID int64, token string, platform *string) (usermodel.DeviceTokenInterface, error) {
	query := `INSERT IGNORE INTO device_tokens (user_id, device_token, platform) VALUES (?, ?, ?)`

	var platformArg interface{}
	if platform != nil {
		platformArg = *platform
	}

	result, err := a.db.GetDB().Exec(query, userID, token, platformArg)
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
		Platform: sql.NullString{String: *platform, Valid: platform != nil},
	}

	return devicetokenconverters.DeviceTokenEntityToDomain(entity), nil
}
