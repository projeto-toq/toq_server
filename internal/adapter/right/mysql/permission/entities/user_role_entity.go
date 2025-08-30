package permissionentities

import (
	"time"
)

type UserRoleEntity struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	RoleID    int64      `db:"role_id"`
	IsActive  bool       `db:"is_active"`
	Status    int64      `db:"status"`
	ExpiresAt *time.Time `db:"expires_at"`
}
