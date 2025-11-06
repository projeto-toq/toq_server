package converters

import (
	devicetokenentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/device_token/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// DeviceTokenEntityToDomain converts entity to domain
func DeviceTokenEntityToDomain(entity devicetokenentity.DeviceTokenEntity) usermodel.DeviceTokenInterface {
	dt := usermodel.NewDeviceToken()
	dt.SetID(entity.ID)
	dt.SetUserID(entity.UserID)
	dt.SetDeviceToken(entity.Token)
	dt.SetDeviceID(entity.DeviceID)
	return dt
}

// DeviceTokenEntitiesToDomain converts slice of entities to domain
func DeviceTokenEntitiesToDomain(entities []devicetokenentity.DeviceTokenEntity) []usermodel.DeviceTokenInterface {
	result := make([]usermodel.DeviceTokenInterface, len(entities))
	for i, entity := range entities {
		result[i] = DeviceTokenEntityToDomain(entity)
	}
	return result
}
