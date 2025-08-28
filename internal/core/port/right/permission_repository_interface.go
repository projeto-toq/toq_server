package right

import (
	"context"
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// PermissionRepositoryInterface define as operações de persistência para permissões
type PermissionRepositoryInterface interface {
	// Operações de Role
	CreateRole(ctx context.Context, role permissionmodel.RoleInterface) error
	GetRoleByID(ctx context.Context, roleID int64) (permissionmodel.RoleInterface, error)
	GetRoleBySlug(ctx context.Context, slug string) (permissionmodel.RoleInterface, error)
	UpdateRole(ctx context.Context, role permissionmodel.RoleInterface) error
	DeleteRole(ctx context.Context, roleID int64) error
	ListRoles(ctx context.Context, systemRolesOnly bool) ([]permissionmodel.RoleInterface, error)

	// Operações de Permission
	CreatePermission(ctx context.Context, permission permissionmodel.PermissionInterface) error
	GetPermissionByID(ctx context.Context, permissionID int64) (permissionmodel.PermissionInterface, error)
	GetPermissionsByResource(ctx context.Context, resource string) ([]permissionmodel.PermissionInterface, error)
	GetPermissionByResourceAndAction(ctx context.Context, resource, action string) (permissionmodel.PermissionInterface, error)
	UpdatePermission(ctx context.Context, permission permissionmodel.PermissionInterface) error
	DeletePermission(ctx context.Context, permissionID int64) error
	ListPermissions(ctx context.Context) ([]permissionmodel.PermissionInterface, error)

	// Operações de RolePermission
	CreateRolePermission(ctx context.Context, rolePermission permissionmodel.RolePermissionInterface) error
	GetRolePermissions(ctx context.Context, roleID int64) ([]permissionmodel.RolePermissionInterface, error)
	GetPermissionRoles(ctx context.Context, permissionID int64) ([]permissionmodel.RolePermissionInterface, error)
	UpdateRolePermissionConditions(ctx context.Context, roleID, permissionID int64, conditions map[string]interface{}) error
	DeleteRolePermission(ctx context.Context, roleID, permissionID int64) error

	// Operações de UserRole
	CreateUserRole(ctx context.Context, userRole permissionmodel.UserRoleInterface) error
	GetUserRoles(ctx context.Context, userID int64) ([]permissionmodel.UserRoleInterface, error)
	GetActiveUserRoles(ctx context.Context, userID int64) ([]permissionmodel.UserRoleInterface, error)
	UpdateUserRoleExpiration(ctx context.Context, userID, roleID int64, expiresAt *time.Time) error
	DeleteUserRole(ctx context.Context, userID, roleID int64) error

	// Operações otimizadas
	GetUserPermissionsWithRoles(ctx context.Context, userID int64) ([]permissionmodel.PermissionInterface, []permissionmodel.RoleInterface, error)
	GetUserEffectivePermissions(ctx context.Context, userID int64) ([]permissionmodel.RolePermissionInterface, error)

	// Operações de cache/batch
	GetMultipleRoles(ctx context.Context, roleIDs []int64) ([]permissionmodel.RoleInterface, error)
	GetMultiplePermissions(ctx context.Context, permissionIDs []int64) ([]permissionmodel.PermissionInterface, error)
}
