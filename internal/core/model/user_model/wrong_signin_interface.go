package usermodel

import "time"

type WrongSigninInterface interface {
	GetUserID() int64
	SetUserID(int64)
	GetFailedAttempts() int64
	SetFailedAttempts(int64)
	GetLastAttemptAt() time.Time
	SetLastAttemptAt(time.Time)
}

func NewWrongSignin() WrongSigninInterface {
	return &wrongSigin{}
}
