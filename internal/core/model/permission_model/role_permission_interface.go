package permissionmodel

type RolePermissionInterface interface {
	GetID() int64
	SetID(id int64)
	GetRoleID() int64
	SetRoleID(roleID int64)
	GetPermissionID() int64
	SetPermissionID(permissionID int64)
	GetGranted() bool
	SetGranted(granted bool)
	GetRole() RoleInterface
	SetRole(role RoleInterface)
	GetPermission() PermissionInterface
	SetPermission(permission PermissionInterface)
}
