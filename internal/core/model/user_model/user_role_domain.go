package usermodel

import (
	"time"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

// userRole representa a associação entre usuário e role
type userRole struct {
	id           int64
	userID       int64
	roleID       int64
	isActive     bool
	status       globalmodel.UserRoleStatus
	expiresAt    *time.Time
	blockedUntil *time.Time
	role         permissionmodel.RoleInterface
}

func NewUserRole() UserRoleInterface {
	return &userRole{
		isActive: true,
		status:   globalmodel.StatusActive,
	}
}

func (ur *userRole) GetID() int64 {
	return ur.id
}

func (ur *userRole) SetID(id int64) {
	ur.id = id
}

func (ur *userRole) GetUserID() int64 {
	return ur.userID
}

func (ur *userRole) SetUserID(userID int64) {
	ur.userID = userID
}

func (ur *userRole) GetRoleID() int64 {
	return ur.roleID
}

func (ur *userRole) SetRoleID(roleID int64) {
	ur.roleID = roleID
}

func (ur *userRole) GetIsActive() bool {
	return ur.isActive
}

func (ur *userRole) SetIsActive(isActive bool) {
	ur.isActive = isActive
}

func (ur *userRole) GetStatus() globalmodel.UserRoleStatus {
	return ur.status
}

func (ur *userRole) SetStatus(status globalmodel.UserRoleStatus) {
	ur.status = status
}

func (ur *userRole) GetExpiresAt() *time.Time {
	return ur.expiresAt
}

func (ur *userRole) SetExpiresAt(expiresAt *time.Time) {
	ur.expiresAt = expiresAt
}

func (ur *userRole) GetBlockedUntil() *time.Time {
	return ur.blockedUntil
}

func (ur *userRole) SetBlockedUntil(blockedUntil *time.Time) {
	ur.blockedUntil = blockedUntil
}

func (ur *userRole) GetRole() permissionmodel.RoleInterface {
	return ur.role
}

func (ur *userRole) SetRole(role permissionmodel.RoleInterface) {
	ur.role = role
}
