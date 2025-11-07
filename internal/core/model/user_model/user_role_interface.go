package usermodel

import (
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

type UserRoleInterface interface {
	GetID() int64
	SetID(id int64)
	GetUserID() int64
	SetUserID(userID int64)
	GetRoleID() int64
	SetRoleID(roleID int64)
	GetIsActive() bool
	SetIsActive(isActive bool)
	GetStatus() globalmodel.UserRoleStatus
	SetStatus(status globalmodel.UserRoleStatus)
	GetExpiresAt() *time.Time
	SetExpiresAt(expiresAt *time.Time)
	GetBlockedUntil() *time.Time
	SetBlockedUntil(blockedUntil *time.Time)
	GetRole() permissionmodel.RoleInterface
	SetRole(role permissionmodel.RoleInterface)
}
