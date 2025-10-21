package permissionentities

type RolePermissionEntity struct {
	ID           int64 `db:"id"`
	RoleID       int64 `db:"role_id"`
	PermissionID int64 `db:"permission_id"`
	Granted      bool  `db:"granted"`
}
