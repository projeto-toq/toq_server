package usermodel

type UserRoleInterface interface {
	GetID() int64
	SetID(int64)
	GetUserID() int64
	SetUserID(int64)
	GetBaseRoleID() int64
	SetBaseRoleID(int64)
	GetRole() UserRole
	SetRole(UserRole)
	IsActive() bool
	SetActive(bool)
	GetStatus() UserRoleStatus
	SetStatus(UserRoleStatus)
	GetStatusReason() string
	SetStatusReason(string)
}

func NewUserRole() UserRoleInterface {
	return &userRole{}
}
