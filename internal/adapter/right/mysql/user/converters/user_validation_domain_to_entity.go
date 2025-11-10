package userconverters

import (
	"database/sql"

	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// UserValidationDomainToEntity converts a domain model to a database entity
//
// This converter handles the translation from clean domain types to database-specific
// types (sql.Null*) for the temp_user_validations table.
//
// Conversion Rules:
//   - string → sql.NullString (Valid=true if non-empty)
//   - time.Time → sql.NullTime (Valid=true if not zero time)
//   - All fields are nullable to support partial validation states
//
// Parameters:
//   - domain: ValidationInterface from core layer
//
// Returns:
//   - entity: UserValidationEntity ready for database operations
//
// Important:
//   - UserID is always required (primary key)
//   - Empty strings are converted to NULL for optional fields
//   - Zero times (IsZero()) are converted to NULL for expiration fields
//   - Codes should be hashed BEFORE calling this converter (service responsibility)
//
// Example:
//
//	validation := usermodel.NewValidation()
//	validation.SetUserID(123)
//	validation.SetNewEmail("newemail@example.com")
//	validation.SetEmailCode("hashed_code_here")
//	validation.SetEmailCodeExp(time.Now().Add(5*time.Minute))
//	entity := UserValidationDomainToEntity(validation)
func UserValidationDomainToEntity(domain usermodel.ValidationInterface) (entity userentity.UserValidationEntity) {
	entity = userentity.UserValidationEntity{}
	entity.UserID = uint32(domain.GetUserID())
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
