package permissionservice

import (
	"context"
	"database/sql"

	cacheport "github.com/projeto-toq/toq_server/internal/core/cache"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	permissionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/permission_repository"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
)

type permissionServiceImpl struct {
	permissionRepository permissionrepository.PermissionRepositoryInterface
	cache                cacheport.CacheInterface
	globalService        globalservice.GlobalServiceInterface
	metrics              metricsport.MetricsPortInterface
}

func NewPermissionService(
	pr permissionrepository.PermissionRepositoryInterface,
	cache cacheport.CacheInterface,
	gs globalservice.GlobalServiceInterface,
	metrics metricsport.MetricsPortInterface,
) PermissionServiceInterface {
	return &permissionServiceImpl{
		permissionRepository: pr,
		cache:                cache,
		globalService:        gs,
		metrics:              metrics,
	}
}

type PermissionServiceInterface interface {
	// Helper para HTTP
	HasHTTPPermission(ctx context.Context, userID int64, method, path string) (bool, error)

	// Gestão de roles
	CreateRole(ctx context.Context, name string, slug permissionmodel.RoleSlug, description string, isSystemRole bool) (permissionmodel.RoleInterface, error)
	UpdateRole(ctx context.Context, input UpdateRoleInput) (permissionmodel.RoleInterface, error)
	DeleteRole(ctx context.Context, roleID int64) error
	RestoreRole(ctx context.Context, roleID int64) (permissionmodel.RoleInterface, error)
	ListRoles(ctx context.Context, input ListRolesInput) (ListRolesOutput, error)
	GetRoleByID(ctx context.Context, roleID int64) (permissionmodel.RoleInterface, error)
	GetRoleByIDWithTx(ctx context.Context, tx *sql.Tx, roleID int64) (permissionmodel.RoleInterface, error)
	GetRoleBySlug(ctx context.Context, slug permissionmodel.RoleSlug) (permissionmodel.RoleInterface, error)
	GetRoleBySlugWithTx(ctx context.Context, tx *sql.Tx, slug permissionmodel.RoleSlug) (permissionmodel.RoleInterface, error)

	// Gestão de permissões
	ListPermissions(ctx context.Context, input ListPermissionsInput) (ListPermissionsOutput, error)
	GetPermissionByID(ctx context.Context, permissionID int64) (permissionmodel.PermissionInterface, error)
	CreatePermission(ctx context.Context, input CreatePermissionInput) (permissionmodel.PermissionInterface, error)
	UpdatePermission(ctx context.Context, input UpdatePermissionInput) (permissionmodel.PermissionInterface, error)
	DeletePermission(ctx context.Context, permissionID int64) error
	GrantPermissionToRole(ctx context.Context, roleID, permissionID int64) error
	RevokePermissionFromRole(ctx context.Context, roleID, permissionID int64) error
	ListRolePermissions(ctx context.Context, input ListRolePermissionsInput) (ListRolePermissionsOutput, error)
	CreateRolePermission(ctx context.Context, input CreateRolePermissionInput) (permissionmodel.RolePermissionInterface, error)
	UpdateRolePermission(ctx context.Context, input UpdateRolePermissionInput) (permissionmodel.RolePermissionInterface, error)
	DeleteRolePermission(ctx context.Context, rolePermissionID int64) error
	GetRolePermissionByID(ctx context.Context, rolePermissionID int64) (permissionmodel.RolePermissionInterface, error)

	// Cache management
	// InvalidateUserCache invalida o cache de permissões de um usuário
	// Retorna erro se a operação falhar; caller deve decidir como tratar.
	//
	// Parameters:
	//   - ctx: Context for tracing and logging
	//   - userID: ID do usuário cujo cache será invalidado
	//   - source: Identificador da operação que causou a invalidação (ex: "assign_role", "switch_active_role")
	//            Usado para rastreabilidade em logs e métricas
	//
	// Returns:
	//   - error: Infrastructure error se a invalidação falhar
	//
	// Usage:
	//   // Quando a falha DEVE ser propagada (operações críticas)
	//   if err := ps.InvalidateUserCache(ctx, userID, "critical_operation"); err != nil {
	//       logger.Error("cache_invalidation_failed", "error", err)
	//       return err
	//   }
	InvalidateUserCache(ctx context.Context, userID int64, source string) error

	// InvalidateUserCacheSafe tenta invalidar o cache mas NÃO propaga erros
	// Registra apenas WARN em logs se a operação falhar.
	// Ideal para operações best-effort após commits de transação.
	//
	// Parameters:
	//   - ctx: Context for tracing and logging
	//   - userID: ID do usuário cujo cache será invalidado
	//   - source: Identificador da operação que causou a invalidação
	//
	// Usage:
	//   // Após commit bem-sucedido, quando falha não deve bloquear fluxo
	//   ps.InvalidateUserCacheSafe(ctx, userID, "assign_role")
	InvalidateUserCacheSafe(ctx context.Context, userID int64, source string)

	ClearUserPermissionsCache(ctx context.Context, userID int64) error
	RefreshUserPermissions(ctx context.Context, userID int64) error

	// Query helpers

	GetUserPermissions(ctx context.Context, userID int64) ([]permissionmodel.PermissionInterface, error)
	GetRolePermissions(ctx context.Context, roleID int64) ([]permissionmodel.PermissionInterface, error)
	GetUserPermissionsWithTx(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.PermissionInterface, error)
}
