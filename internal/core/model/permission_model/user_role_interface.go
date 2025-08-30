package permissionmodel

import "time"

type UserRoleInterface interface {
	GetID() int64
	SetID(id int64)
	GetUserID() int64
	SetUserID(userID int64)
	GetRoleID() int64
	SetRoleID(roleID int64)
	GetIsActive() bool
	SetIsActive(isActive bool)
	GetStatus() UserRoleStatus
	SetStatus(status UserRoleStatus)
	GetExpiresAt() *time.Time
	SetExpiresAt(expiresAt *time.Time)
	GetRole() RoleInterface
	SetRole(role RoleInterface)
}
