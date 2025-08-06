package userconverters

import (
	"log/slog"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UserValidationEntityToDomain(entity []any) (val usermodel.ValidationInterface, err error) {
	val = usermodel.NewValidation()

	user_id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting user_id to int64", "value", entity[0])
		return nil, status.Error(codes.Internal, "Error converting user_id to int64")
	}
	val.SetUserID(user_id)

	if entity[1] != nil {
		new_email, ok := entity[1].([]byte)
		if !ok {
			slog.Error("Error converting new_email to []byte", "value", entity[1])
			return nil, status.Error(codes.Internal, "Error converting new_email to []byte")
		}
		val.SetNewEmail(string(new_email))
	}

	if entity[2] != nil {
		email_code, ok := entity[2].([]byte)
		if !ok {
			slog.Error("Error converting email_code to []byte", "value", entity[2])
			return nil, status.Error(codes.Internal, "Error converting email_code to []byte")
		}
		val.SetEmailCode(string(email_code))
	}

	if entity[3] != nil {
		email_code_exp, ok := entity[3].(time.Time)
		if !ok {
			slog.Error("Error converting email_code_exp to time.Time", "value", entity[3])
			return nil, status.Error(codes.Internal, "Error converting email_code_exp to time.Time")
		}
		val.SetEmailCodeExp(email_code_exp)
	}

	if entity[4] != nil {
		new_phone, ok := entity[4].([]byte)
		if !ok {
			slog.Error("Error converting new_phone to []byte", "value", entity[4])
			return nil, status.Error(codes.Internal, "Error converting new_phone to []byte")
		}
		val.SetNewPhone(string(new_phone))
	}

	if entity[5] != nil {
		phone_code, ok := entity[5].([]byte)
		if !ok {
			slog.Error("Error converting phone_code to []byte", "value", entity[5])
			return nil, status.Error(codes.Internal, "Error converting phone_code to []byte")
		}
		val.SetPhoneCode(string(phone_code))
	}

	if entity[6] != nil {
		phone_code_exp, ok := entity[6].(time.Time)
		if !ok {
			slog.Error("Error converting phone_code_exp to time.Time", "value", entity[6])
			return nil, status.Error(codes.Internal, "Error converting phone_code_exp to time.Time")
		}
		val.SetPhoneCodeExp(phone_code_exp)
	}

	if entity[7] != nil {
		password_code, ok := entity[7].([]byte)
		if !ok {
			slog.Error("Error converting password_code to []byte", "value", entity[7])
			return nil, status.Error(codes.Internal, "Error converting password_code to []byte")
		}
		val.SetPasswordCode(string(password_code))
	}

	if entity[8] != nil {
		password_code_exp, ok := entity[8].(time.Time)
		if !ok {
			slog.Error("Error converting password_code_exp to time.Time", "value", entity[8])
			return nil, status.Error(codes.Internal, "Error converting password_code_exp to time.Time")
		}
		val.SetPasswordCodeExp(password_code_exp)
	}

	return
}
