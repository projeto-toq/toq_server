package userconverters

import (
	"database/sql"

	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// UserRoleDomainToEntity converte UserRoleInterface para UserRoleEntity
func UserRoleDomainToEntity(userRole usermodel.UserRoleInterface) (*userentity.UserRoleEntity, error) {
	if userRole == nil {
		return nil, nil
	}

	entity := &userentity.UserRoleEntity{
		ID:       userRole.GetID(),
		UserID:   userRole.GetUserID(),
		RoleID:   userRole.GetRoleID(),
		IsActive: userRole.GetIsActive(),
		Status:   int64(userRole.GetStatus()),
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
