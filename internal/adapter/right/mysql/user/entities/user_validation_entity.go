package userentity

import "database/sql"

type UserValidationEntity struct {
	UserID          int64
	NewEmail        sql.NullString
	EmailCode       sql.NullString
	EmailCodeExp    sql.NullTime
	NewPhone        sql.NullString
	PhoneCode       sql.NullString
	PhoneCodeExp    sql.NullTime
	PasswordCode    sql.NullString
	PasswordCodeExp sql.NullTime
}
