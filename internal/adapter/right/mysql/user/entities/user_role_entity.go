package userentity

import "database/sql"

type UserRoleEntity struct {
	ID           int64
	UserID       int64
	BaseRoleID   int64
	Role         uint8
	Active       uint8
	Status       uint8
	StatusReason sql.NullString
}
