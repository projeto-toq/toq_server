package permissionconverters

import (
	permissionentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/permission/entities"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

// UserRoleEntityToDomain converte UserRoleEntity para UserRoleInterface
func UserRoleEntityToDomain(entity *permissionentities.UserRoleEntity) (permissionmodel.UserRoleInterface, error) {
	if entity == nil {
		return nil, nil
	}

	userRole := permissionmodel.NewUserRole()
	userRole.SetID(entity.ID)
	userRole.SetUserID(entity.UserID)
	userRole.SetRoleID(entity.RoleID)
	userRole.SetIsActive(entity.IsActive)
	userRole.SetStatus(permissionmodel.UserRoleStatus(entity.Status))

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
