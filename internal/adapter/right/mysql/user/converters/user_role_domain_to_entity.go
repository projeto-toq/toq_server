package userconverters

import (
	"database/sql"

	userentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func UserRoleDomainToEntity(domain usermodel.UserRoleInterface) (entity userentity.UserRoleEntity) {
	entity = userentity.UserRoleEntity{}
	entity.ID = domain.GetID()
	entity.UserID = domain.GetUserID()
	entity.BaseRoleID = domain.GetBaseRoleID()
	entity.Role = uint8(domain.GetRole())
	if domain.IsActive() {
		entity.Active = 1
	} else {
		entity.Active = 0
	}
	entity.Status = uint8(domain.GetStatus())
	entity.StatusReason = sql.NullString{String: domain.GetStatusReason(), Valid: domain.GetStatusReason() != ""}
	return
}
