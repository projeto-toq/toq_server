package userconverters

import (
	"database/sql"

	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// UserRoleDomainToEntity converts a domain UserRoleInterface to a database entity
//
// This converter handles the translation from clean domain types to database-specific
// types (sql.Null*) for the user_roles table.
//
// Conversion Rules:
//   - *time.Time → sql.NullTime (Valid=true if pointer not nil)
//   - UserRoleStatus (enum) → int64 (cast to underlying type)
//   - nil UserRoleInterface → nil entity (safe handling)
//
// Parameters:
//   - userRole: UserRoleInterface from core layer (can be nil)
//
// Returns:
//   - entity: Pointer to UserRoleEntity ready for database operations (nil if input nil)
//   - error: Always nil (kept for interface compatibility)
//
// Important:
//   - ID may be 0 for new records (populated by AUTO_INCREMENT)
//   - Nil pointers for ExpiresAt/BlockedUntil are converted to NULL in database
//   - Status must be valid UserRoleStatus enum value (0=pending, 1=approved, 2=rejected, 3=suspended)
//   - Service layer is responsible for validating Status before conversion
//
// Example:
//
//	userRole := usermodel.NewUserRole()
//	userRole.SetUserID(123)
//	userRole.SetRoleID(456)
//	userRole.SetStatus(globalmodel.UserRoleStatusApproved)
//	entity, _ := UserRoleDomainToEntity(userRole)
func UserRoleDomainToEntity(userRole usermodel.UserRoleInterface) (*userentity.UserRoleEntity, error) {
	if userRole == nil {
		return nil, nil
	}

	entity := &userentity.UserRoleEntity{
		ID:       uint32(userRole.GetID()),
		UserID:   uint32(userRole.GetUserID()),
		RoleID:   uint32(userRole.GetRoleID()),
		IsActive: userRole.GetIsActive(),
		Status:   int8(userRole.GetStatus()),
	}

	// Map optional ExpiresAt field (*time.Time → sql.NullTime)
	if expiresAt := userRole.GetExpiresAt(); expiresAt != nil {
		entity.ExpiresAt = sql.NullTime{
			Time:  *expiresAt,
			Valid: true,
		}
	}

	// Map optional BlockedUntil field (*time.Time → sql.NullTime)
	if blockedUntil := userRole.GetBlockedUntil(); blockedUntil != nil {
		entity.BlockedUntil = sql.NullTime{
			Time:  *blockedUntil,
			Valid: true,
		}
	}

	return entity, nil
}
