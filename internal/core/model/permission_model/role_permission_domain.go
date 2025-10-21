package permissionmodel

// rolePermission representa a associação entre role e permission
type rolePermission struct {
	id           int64
	roleID       int64
	permissionID int64
	granted      bool
	role         RoleInterface
	permission   PermissionInterface
}

func NewRolePermission() RolePermissionInterface {
	return &rolePermission{
		granted: true,
	}
}

func (rp *rolePermission) GetID() int64 {
	return rp.id
}

func (rp *rolePermission) SetID(id int64) {
	rp.id = id
}

func (rp *rolePermission) GetRoleID() int64 {
	return rp.roleID
}

func (rp *rolePermission) SetRoleID(roleID int64) {
	rp.roleID = roleID
}

func (rp *rolePermission) GetPermissionID() int64 {
	return rp.permissionID
}

func (rp *rolePermission) SetPermissionID(permissionID int64) {
	rp.permissionID = permissionID
}

func (rp *rolePermission) GetGranted() bool {
	return rp.granted
}

func (rp *rolePermission) SetGranted(granted bool) {
	rp.granted = granted
}

func (rp *rolePermission) GetRole() RoleInterface {
	return rp.role
}

func (rp *rolePermission) SetRole(role RoleInterface) {
	rp.role = role
}

func (rp *rolePermission) GetPermission() PermissionInterface {
	return rp.permission
}

func (rp *rolePermission) SetPermission(permission PermissionInterface) {
	rp.permission = permission
}
