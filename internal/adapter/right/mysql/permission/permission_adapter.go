package mysqlpermissionadapter

import (
	mysqladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql"
	permissionrepository "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/permission_repository"
)

type PermissionAdapter struct {
	db *mysqladapter.Database
}

func NewPermissionAdapter(db *mysqladapter.Database) permissionrepository.PermissionRepositoryInterface {
	return &PermissionAdapter{
		db: db,
	}
}
