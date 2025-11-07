package userconverters

import (
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// UserRoleEntityToDomain converte UserRoleEntity para UserRoleInterface
func UserRoleEntityToDomain(entity *userentity.UserRoleEntity) (usermodel.UserRoleInterface, error) {
	if entity == nil {
		return nil, nil
	}

	userRole := usermodel.NewUserRole()
	userRole.SetID(entity.ID)
	userRole.SetUserID(entity.UserID)
	userRole.SetRoleID(entity.RoleID)
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
