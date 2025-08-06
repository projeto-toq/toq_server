package usermodel

type userRole struct {
	id           int64
	userID       int64
	baseRoleID   int64
	role         UserRole
	active       bool
	status       UserRoleStatus
	statusReason string
}

func (u *userRole) GetID() int64 {
	return u.id
}

func (u *userRole) SetID(id int64) {
	u.id = id
}

func (u *userRole) GetUserID() int64 {
	return u.userID
}

func (u *userRole) SetUserID(userID int64) {
	u.userID = userID
}

func (u *userRole) GetBaseRoleID() int64 {
	return u.baseRoleID
}

func (u *userRole) SetBaseRoleID(baseRoleID int64) {
	u.baseRoleID = baseRoleID
}

func (u *userRole) GetRole() UserRole {
	return u.role
}

func (u *userRole) SetRole(role UserRole) {
	u.role = role
}

func (u *userRole) IsActive() bool {
	return u.active
}

func (u *userRole) SetActive(active bool) {
	u.active = active
}

func (u *userRole) GetStatus() UserRoleStatus {
	return u.status
}

func (u *userRole) SetStatus(status UserRoleStatus) {
	u.status = status
}

func (u *userRole) GetStatusReason() string {
	return u.statusReason
}

func (u *userRole) SetStatusReason(statusReason string) {
	u.statusReason = statusReason
}
