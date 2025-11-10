package userconverters

import (
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// UserRoleEntityToDomain converts a database UserRoleEntity to domain UserRoleInterface
//
// This converter handles the translation from database-specific types (sql.Null*)
// to clean domain types, ensuring the core layer remains decoupled from database concerns.
//
// Conversion Rules:
//   - sql.NullTime → *time.Time (nil if not Valid, pointer to Time if Valid)
//   - int64 → UserRoleStatus enum (cast to enum type)
//   - nil entity → nil domain (safe handling)
//
// Parameters:
//   - entity: Pointer to UserRoleEntity from database query (can be nil)
//
// Returns:
//   - userRole: UserRoleInterface with all fields populated from entity (nil if input nil)
//   - error: Always nil (kept for interface compatibility)
//
// Important:
//   - Role field is NOT set here - must be populated by Permission Service
//   - Status is cast to UserRoleStatus without validation (service layer responsibility)
//   - NULL timestamps are converted to nil pointers (not zero time)
//
// Example:
//
//	entity := &userentity.UserRoleEntity{
//	    ID: 123,
//	    UserID: 456,
//	    RoleID: 789,
//	    IsActive: true,
//	    Status: 1, // Approved
//	}
//	userRole, _ := UserRoleEntityToDomain(entity)
//	// userRole.GetStatus() == globalmodel.UserRoleStatusApproved
func UserRoleEntityToDomain(entity *userentity.UserRoleEntity) (usermodel.UserRoleInterface, error) {
	if entity == nil {
		return nil, nil
	}

	userRole := usermodel.NewUserRole()
	userRole.SetID(int64(entity.ID))
	userRole.SetUserID(int64(entity.UserID))
	userRole.SetRoleID(int64(entity.RoleID))
	userRole.SetIsActive(entity.IsActive)
	userRole.SetStatus(globalmodel.UserRoleStatus(entity.Status))

	// Map optional ExpiresAt field (sql.NullTime → *time.Time)
	if entity.ExpiresAt.Valid {
		userRole.SetExpiresAt(&entity.ExpiresAt.Time)
	}

	// Map optional BlockedUntil field (sql.NullTime → *time.Time)
	if entity.BlockedUntil.Valid {
		userRole.SetBlockedUntil(&entity.BlockedUntil.Time)
	}

	return userRole, nil
}
