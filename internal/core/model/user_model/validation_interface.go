package usermodel

import "time"

type ValidationInterface interface {
	GetUserID() int64
	SetUserID(int64)
	GetNewEmail() string
	SetNewEmail(string)
	GetEmailCode() string
	SetEmailCode(string)
	GetEmailCodeExp() time.Time
	SetEmailCodeExp(time.Time)
	GetNewPhone() string
	SetNewPhone(string)
	GetPhoneCode() string
	SetPhoneCode(string)
	GetPhoneCodeExp() time.Time
	SetPhoneCodeExp(time.Time)
	GetPasswordCode() string
	SetPasswordCode(string)
	GetPasswordCodeExp() time.Time
	SetPasswordCodeExp(time.Time)
}

func NewValidation() ValidationInterface {
	return &validation{}
}
