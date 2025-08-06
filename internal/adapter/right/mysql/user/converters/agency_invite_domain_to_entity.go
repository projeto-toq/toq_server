package userconverters

import (
	userentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func AgencyInviteDomainToEntity(domain usermodel.InviteInterface) (entity userentity.AgencyInvite) {
	entity = userentity.AgencyInvite{}
	entity.ID = domain.GetID()
	entity.AgencyID = domain.GetAgencyID()
	entity.PhoneNumber = domain.GetPhoneNumber()

	return
}
