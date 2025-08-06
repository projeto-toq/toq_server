package userconverters

import (
	"database/sql"

	userentity "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func UserDomainToEntity(domain usermodel.UserInterface) (entity userentity.UserEntity) {
	entity = userentity.UserEntity{}
	entity.ID = domain.GetID()
	entity.FullName = domain.GetFullName()
	entity.NickName = sql.NullString{String: domain.GetNickName(), Valid: domain.GetNickName() != ""}
	entity.NationalID = domain.GetNationalID()
	entity.CreciNumber = sql.NullString{String: domain.GetCreciNumber(), Valid: domain.GetCreciNumber() != ""}
	entity.CreciState = sql.NullString{String: domain.GetCreciState(), Valid: domain.GetCreciState() != ""}
	entity.CreciValidity = sql.NullTime{Time: domain.GetCreciValidity(), Valid: !domain.GetCreciValidity().IsZero()}
	entity.BornAT = domain.GetBornAt()
	entity.PhoneNumber = domain.GetPhoneNumber()
	entity.Email = domain.GetEmail()
	entity.ZipCode = domain.GetZipCode()
	entity.Street = domain.GetStreet()
	entity.Number = domain.GetNumber()
	entity.Complement = sql.NullString{String: domain.GetComplement(), Valid: domain.GetComplement() != ""}
	entity.Neighborhood = domain.GetNeighborhood()
	entity.City = domain.GetCity()
	entity.State = domain.GetState()
	entity.Photo = sql.NullString{String: string(domain.GetPhoto()), Valid: domain.GetPhoto() != nil}
	entity.Password = domain.GetPassword()
	entity.DeviceToken = sql.NullString{String: string(domain.GetDeviceToken()), Valid: domain.GetDeviceToken() != ""}
	entity.LastActivityAT = domain.GetLastActivityAt()
	entity.Deleted = domain.IsDeleted()
	entity.LastSignInAttempt = sql.NullTime{Time: domain.GetLastSignInAttempt(), Valid: !domain.GetLastSignInAttempt().IsZero()}

	return
}
