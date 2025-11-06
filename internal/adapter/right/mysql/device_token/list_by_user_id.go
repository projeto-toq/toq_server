package mysqldevicetokenadapter

import (
	devicetokenconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/device_token/converters"
	devicetokenentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/device_token/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// ListByUserID returns all device tokens for a user
func (a *DeviceTokenAdapter) ListByUserID(userID int64) ([]usermodel.DeviceTokenInterface, error) {
	query := `SELECT id, user_id, device_token, device_id, platform FROM device_tokens WHERE user_id = ?`

	rows, err := a.db.GetDB().Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []devicetokenentity.DeviceTokenEntity
	for rows.Next() {
		var entity devicetokenentity.DeviceTokenEntity
		if err := rows.Scan(&entity.ID, &entity.UserID, &entity.Token, &entity.DeviceID, &entity.Platform); err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return devicetokenconverters.DeviceTokenEntitiesToDomain(entities), nil
}
