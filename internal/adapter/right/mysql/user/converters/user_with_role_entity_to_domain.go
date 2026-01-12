package userconverters

import (
	"database/sql"
	"fmt"

	permissionconverters "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/converters"
	permissionentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// UserWithRoleEntityToDomain converts a JOIN query result to a complete domain user with active role.
//
// This converter handles the complex mapping from a denormalized database result (users + user_roles + roles)
// to clean domain models (UserInterface with nested UserRoleInterface and RoleInterface).
//
// Conversion Logic:
//  1. Convert user fields (always present) → UserInterface
//  2. Check if user_role fields are present (UserRoleID.Valid)
//  3. If present: convert user_role fields → UserRoleInterface
//  4. Check if role fields are present (RoleID.Valid)
//  5. If present: convert role fields → RoleInterface and attach to UserRoleInterface
//  6. Attach UserRoleInterface to UserInterface.SetActiveRole()
//
// NULL Handling:
//   - User fields: Uses existing UserEntityToDomain logic (sql.Null* → domain types)
//   - UserRole fields: Only converted if UserRoleID.Valid == true
//   - Role fields: Only converted if RoleID.Valid == true
//
// Parameters:
//   - entity: UserWithRoleEntity from JOIN query result
//
// Returns:
//   - user: UserInterface with all fields populated, including ActiveRole if present
//   - error: Conversion errors (should not happen with valid schema)
//
// Example:
//
//	entity := UserWithRoleEntity{...} // from scanUserWithRoleEntity()
//	user, err := UserWithRoleEntityToDomain(entity)
//	if err != nil {
//	    return nil, fmt.Errorf("convert entity: %w", err)
//	}
//	// user.GetActiveRole() is populated if user has active role
func UserWithRoleEntityToDomain(entity userentity.UserWithRoleEntity) (usermodel.UserInterface, error) {
	// Step 1: Convert user fields (always present)
	userEntity := userentity.UserEntity{
		ID:                 entity.UserID,
		FullName:           entity.FullName,
		NickName:           entity.NickName,
		NationalID:         entity.NationalID,
		CreciNumber:        entity.CreciNumber,
		CreciState:         entity.CreciState,
		CreciValidity:      entity.CreciValidity,
		BornAt:             entity.BornAt,
		PhoneNumber:        entity.PhoneNumber,
		Email:              entity.Email,
		ZipCode:            entity.ZipCode,
		Street:             entity.Street,
		Number:             entity.Number,
		Complement:         entity.Complement,
		Neighborhood:       entity.Neighborhood,
		City:               entity.City,
		State:              entity.State,
		Password:           entity.Password,
		OptStatus:          entity.OptStatus,
		LastActivityAt:     entity.LastActivityAt,
		Deleted:            entity.Deleted,
		BlockedUntil:       entity.BlockedUntil,
		PermanentlyBlocked: entity.PermanentlyBlocked,
		CreatedAt:          entity.CreatedAt,
	}

	user := UserEntityToDomain(userEntity)

	// Step 2: Check if user has active role (UserRoleID must be valid)
	if !entity.UserRoleID.Valid {
		// User exists but has no active role (valid state, service layer decides if error)
		return user, nil
	}

	// Step 3: Convert user_role fields
	userRoleEntity := userentity.UserRoleEntity{
		ID:       uint32(entity.UserRoleID.Int32),
		UserID:   uint32(entity.UserRoleUserID.Int32),
		RoleID:   uint32(entity.UserRoleRoleID.Int32),
		IsActive: entity.UserRoleIsActive.Bool,
		Status:   int8(entity.UserRoleStatus.Int16),
	}

	// Handle nullable fields for user_role
	if entity.UserRoleExpiresAt.Valid {
		userRoleEntity.ExpiresAt = sql.NullTime{
			Time:  entity.UserRoleExpiresAt.Time,
			Valid: true,
		}
	}

	userRole, err := UserRoleEntityToDomain(&userRoleEntity)
	if err != nil {
		return nil, fmt.Errorf("convert user_role entity: %w", err)
	}

	// Step 4: Check if role information is present (must be valid if user_role exists)
	if !entity.RoleID.Valid {
		// Inconsistent state: user_role exists but role doesn't (data integrity issue)
		return nil, fmt.Errorf("user_role exists but role is missing (user_id=%d, user_role_id=%d)",
			entity.UserID, entity.UserRoleID.Int32)
	}

	// Step 5: Convert role fields
	roleEntity := permissionentity.RoleEntity{
		ID:           int64(entity.RoleID.Int32),
		Slug:         entity.RoleSlug.String,
		Name:         entity.RoleName.String,
		Description:  entity.RoleDescription.String,
		IsSystemRole: entity.RoleIsSystemRole.Bool,
		IsActive:     entity.RoleIsActive.Bool,
	}

	role := permissionconverters.RoleEntityToDomain(&roleEntity)

	// Step 6: Attach role to user_role and user_role to user
	userRole.SetRole(role)
	user.SetActiveRole(userRole)

	return user, nil
}
