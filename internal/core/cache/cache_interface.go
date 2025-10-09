package cache

import (
	"context"
	"time"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/redis/go-redis/v9"
)

// CacheInterface define a interface para operações de cache Redis
// Focado em cache de permissões modernas usando permission_model
type CacheInterface interface {
	// Métodos de cache de permissões (principal funcionalidade)
	Get(ctx context.Context, userID int64, resource permissionmodel.PermissionResource, action permissionmodel.PermissionAction) (allowed bool, valid bool, err error)
	Set(ctx context.Context, userID int64, resource permissionmodel.PermissionResource, action permissionmodel.PermissionAction, allowed bool, roleSlug permissionmodel.RoleSlug) error

	// Métodos de cache de roles
	GetRole(ctx context.Context, roleSlug permissionmodel.RoleSlug) (roleID int64, valid bool, err error)
	SetRole(ctx context.Context, roleID int64, roleName string, roleSlug permissionmodel.RoleSlug, description string, isActive bool) error

	// Métodos de limpeza específica
	CleanByUser(ctx context.Context, userID int64)                                    // Limpar cache de um usuário específico
	CleanByResource(ctx context.Context, resource permissionmodel.PermissionResource) // Limpar cache de um recurso específico
	Clean(ctx context.Context)                                                        // Limpeza geral do cache Redis

	// Métodos de cache de permissões de usuário
	GetUserPermissions(ctx context.Context, userID int64) ([]byte, error)
	SetUserPermissions(ctx context.Context, userID int64, permissionsJSON []byte, ttl time.Duration) error
	DeleteUserPermissions(ctx context.Context, userID int64) error

	// Métodos de administração
	Close() error // Fechar conexão Redis

	// Injeção de dependências
	SetGlobalService(globalService globalservice.GlobalServiceInterface)

	// Acesso ao cliente Redis para componentes internos
	GetRedisClient() *redis.Client
}
