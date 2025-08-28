package permissionentities

import (
	"time"
)

type UserRoleEntity struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	RoleID    int64      `db:"role_id"`
	IsActive  bool       `db:"is_active"`
	ExpiresAt *time.Time `db:"expires_at"`
}
