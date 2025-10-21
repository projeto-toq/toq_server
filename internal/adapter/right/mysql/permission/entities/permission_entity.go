package permissionentities

// PermissionEntity representa a linha da tabela permissions após a simplificação HTTP-only.
type PermissionEntity struct {
	ID          int64   `db:"id"`
	Name        string  `db:"name"`
	Action      string  `db:"action"`
	Description *string `db:"description"`
	IsActive    bool    `db:"is_active"`
}
