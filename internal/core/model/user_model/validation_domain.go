package usermodel

import (
	"time"
)

type validation struct {
	userID          int64
	newEmail        string
	emailCode       string
	emailCodeExp    time.Time
	newPhone        string
	phoneCode       string
	phoneCodeExp    time.Time
	passwordCode    string
	passwordCodeExp time.Time
}

func (v *validation) GetNewEmail() string {
	return v.newEmail
}

func (v *validation) SetNewEmail(newEmail string) {
	v.newEmail = newEmail
}

func (v *validation) GetEmailCode() string {
	return v.emailCode
}

func (v *validation) SetEmailCode(emailCode string) {
	v.emailCode = emailCode
}

func (v *validation) GetEmailCodeExp() time.Time {
	return v.emailCodeExp
}

func (v *validation) SetEmailCodeExp(emailCodeExp time.Time) {
	v.emailCodeExp = emailCodeExp
}

func (v *validation) GetNewPhone() string {
	return v.newPhone
}

func (v *validation) SetNewPhone(newPhone string) {
	v.newPhone = newPhone
}

func (v *validation) GetPhoneCode() string {
	return v.phoneCode
}

func (v *validation) SetPhoneCode(phoneCode string) {
	v.phoneCode = phoneCode
}

func (v *validation) GetPhoneCodeExp() time.Time {
	return v.phoneCodeExp
}

func (v *validation) SetPhoneCodeExp(phoneCodeExp time.Time) {
	v.phoneCodeExp = phoneCodeExp
}

func (v *validation) GetPasswordCode() string {
	return v.passwordCode
}

func (v *validation) SetPasswordCode(passwordCode string) {
	v.passwordCode = passwordCode
}

func (v *validation) GetPasswordCodeExp() time.Time {
	return v.passwordCodeExp
}

func (v *validation) SetPasswordCodeExp(passwordCodeExp time.Time) {
	v.passwordCodeExp = passwordCodeExp
}

func (v *validation) GetUserID() int64 {
	return v.userID
}

func (v *validation) SetUserID(userID int64) {
	v.userID = userID
}
