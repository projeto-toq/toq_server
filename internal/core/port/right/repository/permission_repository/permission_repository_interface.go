package permissionrepository

import (
	"context"
	"database/sql"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
)

type PermissionRepositoryInterface interface {
	// Role operations
	CreateRole(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleInterface) error
	GetRoleByID(ctx context.Context, tx *sql.Tx, roleID int64) (permissionmodel.RoleInterface, error)
	GetRoleBySlug(ctx context.Context, tx *sql.Tx, slug string) (permissionmodel.RoleInterface, error)
	GetAllRoles(ctx context.Context, tx *sql.Tx) ([]permissionmodel.RoleInterface, error)
	ListRoles(ctx context.Context, tx *sql.Tx, filter RoleListFilter) (RoleListResult, error)
	UpdateRole(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleInterface) error
	DeleteRole(ctx context.Context, tx *sql.Tx, roleID int64) error

	// Permission operations
	CreatePermission(ctx context.Context, tx *sql.Tx, permission permissionmodel.PermissionInterface) error
	GetPermissionByID(ctx context.Context, tx *sql.Tx, permissionID int64) (permissionmodel.PermissionInterface, error)
	GetPermissionByName(ctx context.Context, tx *sql.Tx, name string) (permissionmodel.PermissionInterface, error)
	GetPermissionByAction(ctx context.Context, tx *sql.Tx, action string) (permissionmodel.PermissionInterface, error)
	GetAllPermissions(ctx context.Context, tx *sql.Tx) ([]permissionmodel.PermissionInterface, error)
	ListPermissions(ctx context.Context, tx *sql.Tx, filter PermissionListFilter) (PermissionListResult, error)
	UpdatePermission(ctx context.Context, tx *sql.Tx, permission permissionmodel.PermissionInterface) error
	DeletePermission(ctx context.Context, tx *sql.Tx, permissionID int64) error

	// RolePermission operations
	CreateRolePermission(ctx context.Context, tx *sql.Tx, rolePermission permissionmodel.RolePermissionInterface) error
	GetRolePermissionByID(ctx context.Context, tx *sql.Tx, rolePermissionID int64) (permissionmodel.RolePermissionInterface, error)
	GetRolePermissionsByRoleID(ctx context.Context, tx *sql.Tx, roleID int64) ([]permissionmodel.RolePermissionInterface, error)
	GetGrantedPermissionsByRoleID(ctx context.Context, tx *sql.Tx, roleID int64) ([]permissionmodel.PermissionInterface, error)
	GetRolePermissionByRoleIDAndPermissionID(ctx context.Context, tx *sql.Tx, roleID, permissionID int64) (permissionmodel.RolePermissionInterface, error)
	ListRolePermissions(ctx context.Context, tx *sql.Tx, filter RolePermissionListFilter) (RolePermissionListResult, error)
	UpdateRolePermission(ctx context.Context, tx *sql.Tx, rolePermission permissionmodel.RolePermissionInterface) error
	DeleteRolePermission(ctx context.Context, tx *sql.Tx, rolePermissionID int64) error
	GetRoleIDsByPermissionID(ctx context.Context, tx *sql.Tx, permissionID int64) ([]int64, error)

	// Complex queries for permission checking
	GetUserPermissions(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.PermissionInterface, error)
	GetActiveUserIDsByRoleID(ctx context.Context, tx *sql.Tx, roleID int64) ([]int64, error)
}

type RoleListFilter struct {
	Page         int
	Limit        int
	Name         string
	Slug         string
	Description  string
	IsSystemRole *bool
	IsActive     *bool
	IDFrom       *int64
	IDTo         *int64
}

type RoleListResult struct {
	Roles []permissionmodel.RoleInterface
	Total int64
}

type PermissionListFilter struct {
	Page     int
	Limit    int
	Name     string
	Action   string
	IsActive *bool
}

type PermissionListResult struct {
	Permissions []permissionmodel.PermissionInterface
	Total       int64
}

type RolePermissionListFilter struct {
	Page         int
	Limit        int
	RoleID       *int64
	PermissionID *int64
	Granted      *bool
}

type RolePermissionListResult struct {
	RolePermissions []permissionmodel.RolePermissionInterface
	Total           int64
}
