package userconverters

import (
	userentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func PrivilegeDomainToEntity(domain usermodel.PrivilegeInterface, roleID int64) (entity userentity.PrivilegeEntity) {
	entity = userentity.PrivilegeEntity{}
	entity.ID = domain.ID()
	entity.RoleID = roleID
	entity.Service = uint8(domain.Service())
	entity.Method = uint8(domain.Method())
	if domain.Allowed() {
		entity.Allowed = 1
	} else {
		entity.Allowed = 0
	}

	return
}
