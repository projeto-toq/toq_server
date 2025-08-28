package permissionentities

type PermissionEntity struct {
	ID          int64   `db:"id"`
	Name        string  `db:"name"`
	Slug        string  `db:"slug"`
	Resource    string  `db:"resource"`
	Action      string  `db:"action"`
	Description string  `db:"description"`
	Conditions  *string `db:"conditions"` // JSON string
	IsActive    bool    `db:"is_active"`
}
