package permissionservice

import (
	"context"
	"database/sql"
	"time"

	cacheport "github.com/giulio-alfieri/toq_server/internal/core/cache"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	permissionrepository "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/permission_repository"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
)

type permissionServiceImpl struct {
	permissionRepository permissionrepository.PermissionRepositoryInterface
	cache                cacheport.CacheInterface
	globalService        globalservice.GlobalServiceInterface
}

func NewPermissionService(
	pr permissionrepository.PermissionRepositoryInterface,
	cache cacheport.CacheInterface,
	gs globalservice.GlobalServiceInterface,
) PermissionServiceInterface {
	return &permissionServiceImpl{
		permissionRepository: pr,
		cache:                cache,
		globalService:        gs,
	}
}

type PermissionServiceInterface interface {
	// Verificação principal de permissões
	HasPermission(ctx context.Context, userID int64, resource, action string, permContext *permissionmodel.PermissionContext) (bool, error)

	// Helper para HTTP
	HasHTTPPermission(ctx context.Context, userID int64, method, path string) (bool, error)

	// Gestão de roles
	CreateRole(ctx context.Context, name, slug, description string, isSystemRole bool) (permissionmodel.RoleInterface, error)
	AssignRoleToUser(ctx context.Context, userID, roleID int64, expiresAt *time.Time) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error

	// Gestão de permissões
	GrantPermissionToRole(ctx context.Context, roleID, permissionID int64) error
	RevokePermissionFromRole(ctx context.Context, roleID, permissionID int64) error

	// Cache management
	InvalidateUserCache(ctx context.Context, userID int64) error
	RefreshUserPermissions(ctx context.Context, userID int64) error

	// Query helpers
	GetUserRoles(ctx context.Context, userID int64) ([]permissionmodel.UserRoleInterface, error)
	GetUserPermissions(ctx context.Context, userID int64) ([]permissionmodel.PermissionInterface, error)
	GetRolePermissions(ctx context.Context, roleID int64) ([]permissionmodel.PermissionInterface, error)

	// NOVOS: Controle de múltiplos roles ativos
	SwitchActiveRole(ctx context.Context, userID, newRoleID int64) error
	GetActiveUserRole(ctx context.Context, userID int64) (permissionmodel.UserRoleInterface, error)
	DeactivateAllUserRoles(ctx context.Context, userID int64) error
	GetRoleBySlug(ctx context.Context, slug string) (permissionmodel.RoleInterface, error)

	// Métodos com transação (para uso em fluxos maiores)
	GetRoleBySlugWithTx(ctx context.Context, tx *sql.Tx, slug string) (permissionmodel.RoleInterface, error)
	AssignRoleToUserWithTx(ctx context.Context, tx *sql.Tx, userID, roleID int64, expiresAt *time.Time) error
	GetUserPermissionsWithTx(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.PermissionInterface, error)
	GetUserRolesWithTx(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.UserRoleInterface, error)
	SwitchActiveRoleWithTx(ctx context.Context, tx *sql.Tx, userID, newRoleID int64) error
}
