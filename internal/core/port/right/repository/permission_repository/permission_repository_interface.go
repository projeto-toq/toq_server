package permissionrepository

import (
	"context"
	"database/sql"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

type PermissionRepositoryInterface interface {
	// Role operations
	CreateRole(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleInterface) error
	GetRoleByID(ctx context.Context, tx *sql.Tx, roleID int64) (permissionmodel.RoleInterface, error)
	GetRoleBySlug(ctx context.Context, tx *sql.Tx, slug string) (permissionmodel.RoleInterface, error)
	GetAllRoles(ctx context.Context, tx *sql.Tx) ([]permissionmodel.RoleInterface, error)
	UpdateRole(ctx context.Context, tx *sql.Tx, role permissionmodel.RoleInterface) error
	DeleteRole(ctx context.Context, tx *sql.Tx, roleID int64) error

	// Permission operations
	CreatePermission(ctx context.Context, tx *sql.Tx, permission permissionmodel.PermissionInterface) error
	GetPermissionByID(ctx context.Context, tx *sql.Tx, permissionID int64) (permissionmodel.PermissionInterface, error)
	GetPermissionByName(ctx context.Context, tx *sql.Tx, name string) (permissionmodel.PermissionInterface, error)
	GetPermissionsByResource(ctx context.Context, tx *sql.Tx, resource string) ([]permissionmodel.PermissionInterface, error)
	GetPermissionsByResourceAndAction(ctx context.Context, tx *sql.Tx, resource, action string) ([]permissionmodel.PermissionInterface, error)
	GetAllPermissions(ctx context.Context, tx *sql.Tx) ([]permissionmodel.PermissionInterface, error)
	UpdatePermission(ctx context.Context, tx *sql.Tx, permission permissionmodel.PermissionInterface) error
	DeletePermission(ctx context.Context, tx *sql.Tx, permissionID int64) error

	// UserRole operations
	CreateUserRole(ctx context.Context, tx *sql.Tx, userRole permissionmodel.UserRoleInterface) error
	GetUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.UserRoleInterface, error)
	GetActiveUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.UserRoleInterface, error)
	GetActiveUserRoleByUserID(ctx context.Context, tx *sql.Tx, userID int64) (permissionmodel.UserRoleInterface, error)
	GetUserRoleByUserIDAndRoleID(ctx context.Context, tx *sql.Tx, userID, roleID int64) (permissionmodel.UserRoleInterface, error)
	UpdateUserRole(ctx context.Context, tx *sql.Tx, userRole permissionmodel.UserRoleInterface) error
	DeleteUserRole(ctx context.Context, tx *sql.Tx, userRoleID int64) error
	DeactivateAllUserRoles(ctx context.Context, tx *sql.Tx, userID int64) error
	ActivateUserRole(ctx context.Context, tx *sql.Tx, userID, roleID int64) error

	// RolePermission operations
	CreateRolePermission(ctx context.Context, tx *sql.Tx, rolePermission permissionmodel.RolePermissionInterface) error
	GetRolePermissionsByRoleID(ctx context.Context, tx *sql.Tx, roleID int64) ([]permissionmodel.RolePermissionInterface, error)
	GetGrantedPermissionsByRoleID(ctx context.Context, tx *sql.Tx, roleID int64) ([]permissionmodel.PermissionInterface, error)
	GetRolePermissionByRoleIDAndPermissionID(ctx context.Context, tx *sql.Tx, roleID, permissionID int64) (permissionmodel.RolePermissionInterface, error)
	UpdateRolePermission(ctx context.Context, tx *sql.Tx, rolePermission permissionmodel.RolePermissionInterface) error
	DeleteRolePermission(ctx context.Context, tx *sql.Tx, rolePermissionID int64) error

	// Complex queries for permission checking
	GetUserPermissions(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.PermissionInterface, error)
	HasUserPermission(ctx context.Context, tx *sql.Tx, userID int64, resource, action string) (bool, error)
	GetUserPermissionsByResource(ctx context.Context, tx *sql.Tx, userID int64, resource string) ([]permissionmodel.PermissionInterface, error)
}
