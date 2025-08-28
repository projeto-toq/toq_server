package permissionentities

type RoleEntity struct {
	ID           int64  `db:"id"`
	Name         string `db:"name"`
	Slug         string `db:"slug"`
	Description  string `db:"description"`
	IsSystemRole bool   `db:"is_system_role"`
	IsActive     bool   `db:"is_active"`
}
