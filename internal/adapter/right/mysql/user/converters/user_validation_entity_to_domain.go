package userconverters

import (
	userentity "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/user/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// UserValidationEntityToDomainTyped converts a strongly-typed UserValidationEntity to domain model
//
// This is the preferred conversion function providing type safety and eliminating runtime panics
// from type assertions. Use this function for all new code.
//
// Parameters:
//   - entity: UserValidationEntity from database query
//
// Returns:
//   - val: ValidationInterface with all fields populated from entity
//
// Note: sql.Null* fields are checked for Valid before accessing
func UserValidationEntityToDomainTyped(entity userentity.UserValidationEntity) usermodel.ValidationInterface {
	val := usermodel.NewValidation()

	// Set user ID (always present, primary key)
	val.SetUserID(int64(entity.UserID))

	// Map optional fields - check Valid before accessing
	if entity.NewEmail.Valid {
		val.SetNewEmail(entity.NewEmail.String)
	}

	if entity.EmailCode.Valid {
		val.SetEmailCode(entity.EmailCode.String)
	}

	if entity.EmailCodeExp.Valid {
		val.SetEmailCodeExp(entity.EmailCodeExp.Time)
	}

	if entity.NewPhone.Valid {
		val.SetNewPhone(entity.NewPhone.String)
	}

	if entity.PhoneCode.Valid {
		val.SetPhoneCode(entity.PhoneCode.String)
	}

	if entity.PhoneCodeExp.Valid {
		val.SetPhoneCodeExp(entity.PhoneCodeExp.Time)
	}

	if entity.PasswordCode.Valid {
		val.SetPasswordCode(entity.PasswordCode.String)
	}

	if entity.PasswordCodeExp.Valid {
		val.SetPasswordCodeExp(entity.PasswordCodeExp.Time)
	}

	return val
}
