package userconverters

import (
	"errors"
	"time"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

func UserValidationEntityToDomain(entity []any) (val usermodel.ValidationInterface, err error) {
	val = usermodel.NewValidation()

	user_id, ok := entity[0].(int64)
	if !ok {
		return nil, errors.New("invalid user_id type")
	}
	val.SetUserID(user_id)

	if entity[1] != nil {
		new_email, ok := entity[1].([]byte)
		if !ok {
			return nil, errors.New("invalid new_email type")
		}
		val.SetNewEmail(string(new_email))
	}

	if entity[2] != nil {
		email_code, ok := entity[2].([]byte)
		if !ok {
			return nil, errors.New("invalid email_code type")
		}
		val.SetEmailCode(string(email_code))
	}

	if entity[3] != nil {
		email_code_exp, ok := entity[3].(time.Time)
		if !ok {
			return nil, errors.New("invalid email_code_exp type")
		}
		val.SetEmailCodeExp(email_code_exp)
	}

	if entity[4] != nil {
		new_phone, ok := entity[4].([]byte)
		if !ok {
			return nil, errors.New("invalid new_phone type")
		}
		val.SetNewPhone(string(new_phone))
	}

	if entity[5] != nil {
		phone_code, ok := entity[5].([]byte)
		if !ok {
			return nil, errors.New("invalid phone_code type")
		}
		val.SetPhoneCode(string(phone_code))
	}

	if entity[6] != nil {
		phone_code_exp, ok := entity[6].(time.Time)
		if !ok {
			return nil, errors.New("invalid phone_code_exp type")
		}
		val.SetPhoneCodeExp(phone_code_exp)
	}

	if entity[7] != nil {
		password_code, ok := entity[7].([]byte)
		if !ok {
			return nil, errors.New("invalid password_code type")
		}
		val.SetPasswordCode(string(password_code))
	}

	if entity[8] != nil {
		password_code_exp, ok := entity[8].(time.Time)
		if !ok {
			return nil, errors.New("invalid password_code_exp type")
		}
		val.SetPasswordCodeExp(password_code_exp)
	}

	return
}
