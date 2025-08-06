package userconverters

import (
	userentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func BaseRoleDomainToEntity(domain usermodel.BaseRoleInterface) (entity userentity.BaseRoleEntity) {
	entity = userentity.BaseRoleEntity{}
	entity.ID = domain.GetID()
	entity.Role = uint8(domain.GetRole())
	entity.Name = domain.GetName()

	return
}
