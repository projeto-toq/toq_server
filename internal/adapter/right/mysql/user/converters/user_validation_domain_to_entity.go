package userconverters

import (
	"database/sql"

	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

func UserValidationDomainToEntity(domain usermodel.ValidationInterface) (entity userentity.UserValidationEntity) {
	entity = userentity.UserValidationEntity{}
	entity.UserID = domain.GetUserID()
	entity.NewEmail = sql.NullString{String: domain.GetNewEmail(), Valid: domain.GetNewEmail() != ""}
	entity.EmailCode = sql.NullString{String: domain.GetEmailCode(), Valid: domain.GetEmailCode() != ""}
	entity.EmailCodeExp = sql.NullTime{Time: domain.GetEmailCodeExp(), Valid: !domain.GetEmailCodeExp().IsZero()}
	entity.NewPhone = sql.NullString{String: domain.GetNewPhone(), Valid: domain.GetNewPhone() != ""}
	entity.PhoneCode = sql.NullString{String: domain.GetPhoneCode(), Valid: domain.GetPhoneCode() != ""}
	entity.PhoneCodeExp = sql.NullTime{Time: domain.GetPhoneCodeExp(), Valid: !domain.GetPhoneCodeExp().IsZero()}
	entity.PasswordCode = sql.NullString{String: domain.GetPasswordCode(), Valid: domain.GetPasswordCode() != ""}
	entity.PasswordCodeExp = sql.NullTime{Time: domain.GetPasswordCodeExp(), Valid: !domain.GetPasswordCodeExp().IsZero()}
	return
}
