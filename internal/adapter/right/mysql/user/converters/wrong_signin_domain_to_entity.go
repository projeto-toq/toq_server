package userconverters

import (
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

func WrongSignInDomainToEntity(domain usermodel.WrongSigninInterface) (entity userentity.WrongSignInEntity) {
	entity = userentity.WrongSignInEntity{}
	entity.UserID = domain.GetUserID()
	entity.FailedAttempts = uint8(domain.GetFailedAttempts())
	entity.LastAttemptAT = domain.GetLastAttemptAt()
	return
}
