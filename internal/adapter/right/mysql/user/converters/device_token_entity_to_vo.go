package userconverters

import (
	"database/sql"

	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// DeviceTokenEntityToVO converts a database entity to a domain Value Object
//
// This converter handles the translation from database-specific types (sql.Null*)
// to clean domain types, ensuring the core layer remains decoupled from database concerns.
//
// Conversion Rules:
//   - sql.NullString → string (empty string if NULL)
//   - All fields are mandatory except Platform (nullable in schema)
//
// Parameters:
//   - entity: DeviceTokenEntity from database query
//
// Returns:
//   - token: DeviceToken Value Object with all fields populated
//
// Example:
//
//	entity := DeviceTokenEntity{ID: 1, UserID: 123, Token: "fcm_abc", DeviceID: "uuid", Platform: sql.NullString{String: "android", Valid: true}}
//	vo := DeviceTokenEntityToVO(entity)
//	// vo.Platform == "android"
func DeviceTokenEntityToVO(entity userentity.DeviceTokenEntity) usermodel.DeviceToken {
	vo := usermodel.DeviceToken{
		ID:       entity.ID,
		UserID:   entity.UserID,
		Token:    entity.Token,
		DeviceID: entity.DeviceID,
	}

	// Handle NULL platform (convert to empty string if NULL)
	if entity.Platform.Valid {
		vo.Platform = entity.Platform.String
	}

	return vo
}

// DeviceTokenEntitiesToVOs converts a slice of entities to domain Value Objects
//
// Convenience function for batch conversion (e.g., ListDeviceTokensByUserID).
//
// Parameters:
//   - entities: Slice of DeviceTokenEntity from database query
//
// Returns:
//   - tokens: Slice of DeviceToken Value Objects
func DeviceTokenEntitiesToVOs(entities []userentity.DeviceTokenEntity) []usermodel.DeviceToken {
	vos := make([]usermodel.DeviceToken, len(entities))
	for i, entity := range entities {
		vos[i] = DeviceTokenEntityToVO(entity)
	}
	return vos
}

// DeviceTokenVOToEntity converts a domain Value Object to a database entity
//
// This converter handles the translation from clean domain types to database-specific
// types (sql.Null*), preparing data for database insertion/update.
//
// Conversion Rules:
//   - string → sql.NullString (Valid=true if non-empty)
//   - Empty Platform → sql.NullString{Valid: false} (NULL in DB)
//
// Parameters:
//   - vo: DeviceToken Value Object from domain layer
//
// Returns:
//   - entity: DeviceTokenEntity ready for database operations
//
// Important:
//   - ID may be 0 for new records (populated by AUTO_INCREMENT)
//   - Empty Platform results in NULL (not empty string)
func DeviceTokenVOToEntity(vo usermodel.DeviceToken) userentity.DeviceTokenEntity {
	entity := userentity.DeviceTokenEntity{
		ID:       vo.ID,
		UserID:   vo.UserID,
		Token:    vo.Token,
		DeviceID: vo.DeviceID,
		Platform: sql.NullString{
			String: vo.Platform,
			Valid:  vo.Platform != "",
		},
	}

	return entity
}
