package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/redis/go-redis/v9"
)

// RedisCache implementa CacheInterface usando Redis
type RedisCache struct {
	client        *redis.Client
	globalService globalservice.GlobalServiceInterface
	defaultTTL    time.Duration
	keyPrefix     string
}

// PermissionCacheEntry representa uma entrada de cache para o sistema de permissões moderno
type PermissionCacheEntry struct {
	Allowed   bool                               `json:"allowed"`
	Valid     bool                               `json:"valid"`
	CreatedAt time.Time                          `json:"created_at"`
	UserID    int64                              `json:"user_id"`
	Resource  permissionmodel.PermissionResource `json:"resource"`
	Action    permissionmodel.PermissionAction   `json:"action"`
	RoleSlug  permissionmodel.RoleSlug           `json:"role_slug"`
	ExpiresAt *time.Time                         `json:"expires_at,omitempty"`
}

// RoleCacheEntry representa uma entrada de cache para roles
type RoleCacheEntry struct {
	RoleID      int64                    `json:"role_id"`
	RoleName    string                   `json:"role_name"`
	RoleSlug    permissionmodel.RoleSlug `json:"role_slug"`
	Description string                   `json:"description"`
	IsActive    bool                     `json:"is_active"`
	CachedAt    time.Time                `json:"cached_at"`
}

// NewRedisCache cria uma nova instância do cache Redis
func NewRedisCache(redisURL string, globalService globalservice.GlobalServiceInterface) (CacheInterface, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Testar conexão
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	slog.Info("Redis cache connected successfully", "url", redisURL)

	return &RedisCache{
		client:        client,
		globalService: globalService,
		defaultTTL:    15 * time.Minute, // TTL padrão de 15 minutos
		keyPrefix:     "toq_cache:",
	}, nil
}

// SetGlobalService injeta o GlobalService após a criação do cache
// Usado para resolver dependências circulares entre Cache e GlobalService
func (rc *RedisCache) SetGlobalService(globalService globalservice.GlobalServiceInterface) {
	rc.globalService = globalService
	slog.Debug("GlobalService injected into RedisCache")
}

// GetRedisClient retorna o cliente Redis para uso interno
// Usado pelo ActivityTracker para operações diretas no Redis
func (rc *RedisCache) GetRedisClient() *redis.Client {
	return rc.client
}

// Get recupera uma entrada do cache para o sistema de permissões moderno
func (rc *RedisCache) Get(ctx context.Context, userID int64, resource permissionmodel.PermissionResource, action permissionmodel.PermissionAction) (allowed bool, valid bool, err error) {
	key := rc.buildPermissionKey(userID, resource, action)

	result, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Cache miss - não é erro
			return false, false, nil
		}
		slog.Error("Redis cache get error", "key", key, "error", err)
		return false, false, err
	}

	var entry PermissionCacheEntry
	if err := json.Unmarshal([]byte(result), &entry); err != nil {
		slog.Error("Failed to unmarshal cache entry", "key", key, "error", err)
		return false, false, err
	}

	// Verificar se o cache não expirou
	if time.Since(entry.CreatedAt) > rc.defaultTTL {
		// Cache expirado
		rc.client.Del(ctx, key) // Remover entrada expirada
		return false, false, nil
	}

	slog.Debug("Cache hit", "key", key, "allowed", entry.Allowed)
	return entry.Allowed, entry.Valid, nil
}

// Set armazena uma entrada no cache para o sistema de permissões moderno
func (rc *RedisCache) Set(ctx context.Context, userID int64, resource permissionmodel.PermissionResource, action permissionmodel.PermissionAction, allowed bool, roleSlug permissionmodel.RoleSlug) error {
	key := rc.buildPermissionKey(userID, resource, action)

	entry := PermissionCacheEntry{
		Allowed:   allowed,
		Valid:     true,
		CreatedAt: time.Now(),
		UserID:    userID,
		Resource:  resource,
		Action:    action,
		RoleSlug:  roleSlug,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		slog.Error("Failed to marshal cache entry", "key", key, "error", err)
		return err
	}

	err = rc.client.Set(ctx, key, data, rc.defaultTTL).Err()
	if err != nil {
		slog.Error("Failed to set cache entry", "key", key, "error", err)
		return err
	}

	slog.Debug("Cache set", "key", key, "allowed", allowed)
	return nil
}

// GetRole recupera informações de role do cache
func (rc *RedisCache) GetRole(ctx context.Context, roleSlug permissionmodel.RoleSlug) (roleID int64, valid bool, err error) {
	key := rc.buildRoleKey(roleSlug)

	result, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, false, nil
		}
		slog.Error("Redis cache get role error", "key", key, "error", err)
		return 0, false, err
	}

	var entry RoleCacheEntry
	if err := json.Unmarshal([]byte(result), &entry); err != nil {
		slog.Error("Failed to unmarshal role cache entry", "key", key, "error", err)
		return 0, false, err
	}

	// Verificar se o cache não expirou
	if time.Since(entry.CachedAt) > rc.defaultTTL {
		rc.client.Del(ctx, key)
		return 0, false, nil
	}

	return entry.RoleID, true, nil
}

// SetRole armazena informações de role no cache
func (rc *RedisCache) SetRole(ctx context.Context, roleID int64, roleName string, roleSlug permissionmodel.RoleSlug, description string, isActive bool) error {
	key := rc.buildRoleKey(roleSlug)

	entry := RoleCacheEntry{
		RoleID:      roleID,
		RoleName:    roleName,
		RoleSlug:    roleSlug,
		Description: description,
		IsActive:    isActive,
		CachedAt:    time.Now(),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		slog.Error("Failed to marshal role cache entry", "key", key, "error", err)
		return err
	}

	err = rc.client.Set(ctx, key, data, rc.defaultTTL).Err()
	if err != nil {
		slog.Error("Failed to set role cache entry", "key", key, "error", err)
		return err
	}

	slog.Debug("Role cache set", "key", key, "role_slug", roleSlug)
	return nil
}

// Clean limpa o cache (pode ser específico ou geral)
func (rc *RedisCache) Clean(ctx context.Context) {
	pattern := rc.keyPrefix + "*"

	slog.Debug("Cleaning cache", "pattern", pattern)

	// Usar SCAN para evitar bloquear o Redis
	iter := rc.client.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		slog.Error("Error scanning cache keys", "error", err)
		return
	}

	if len(keys) > 0 {
		deleted, err := rc.client.Del(ctx, keys...).Result()
		if err != nil {
			slog.Error("Error deleting cache keys", "error", err)
			return
		}
		slog.Info("Cache cleaned", "deleted_keys", deleted)
	} else {
		slog.Debug("No cache keys to clean")
	}
}

// CleanByMethod limpa cache específico por recurso e ação
func (rc *RedisCache) CleanByResource(ctx context.Context, resource permissionmodel.PermissionResource) {
	pattern := rc.keyPrefix + fmt.Sprintf("permission:*:resource:%s:*", resource)

	slog.Debug("Cleaning cache by resource", "resource", resource, "pattern", pattern)

	iter := rc.client.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		slog.Error("Error scanning cache keys by resource", "resource", resource, "error", err)
		return
	}

	if len(keys) > 0 {
		deleted, err := rc.client.Del(ctx, keys...).Result()
		if err != nil {
			slog.Error("Error deleting cache keys by resource", "resource", resource, "error", err)
			return
		}
		slog.Debug("Cache cleaned by resource", "resource", resource, "deleted_keys", deleted)
	}
}

// CleanByUser limpa cache específico de um usuário
func (rc *RedisCache) CleanByUser(ctx context.Context, userID int64) {
	pattern := rc.keyPrefix + fmt.Sprintf("permission:user:%d:*", userID)

	slog.Debug("Cleaning cache by user", "user_id", userID, "pattern", pattern)

	iter := rc.client.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		slog.Error("Error scanning cache keys by user", "user_id", userID, "error", err)
		return
	}

	if len(keys) > 0 {
		deleted, err := rc.client.Del(ctx, keys...).Result()
		if err != nil {
			slog.Error("Error deleting cache keys by user", "user_id", userID, "error", err)
			return
		}
		slog.Debug("Cache cleaned by user", "user_id", userID, "deleted_keys", deleted)
	}
}

// buildPermissionKey constrói a chave do cache para permissões
func (rc *RedisCache) buildPermissionKey(userID int64, resource permissionmodel.PermissionResource, action permissionmodel.PermissionAction) string {
	return fmt.Sprintf("%spermission:user:%d:resource:%s:action:%s", rc.keyPrefix, userID, resource, action)
}

// buildRoleKey constrói a chave do cache para roles
func (rc *RedisCache) buildRoleKey(roleSlug permissionmodel.RoleSlug) string {
	return fmt.Sprintf("%srole:slug:%s", rc.keyPrefix, roleSlug)
}

// Close fecha a conexão com o Redis
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}

// GetUserPermissions busca permissões de usuário do Redis
func (rc *RedisCache) GetUserPermissions(ctx context.Context, userID int64) ([]byte, error) {
	key := fmt.Sprintf("%suser_permissions:%d", rc.keyPrefix, userID)

	result, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("cache miss - user permissions not found for user %d", userID)
		}
		return nil, fmt.Errorf("failed to get user permissions from Redis: %w", err)
	}

	slog.Debug("User permissions cache hit", "userID", userID, "dataSize", len(result))
	return []byte(result), nil
}

// SetUserPermissions armazena permissões de usuário no Redis
func (rc *RedisCache) SetUserPermissions(ctx context.Context, userID int64, permissionsJSON []byte, ttl time.Duration) error {
	key := fmt.Sprintf("%suser_permissions:%d", rc.keyPrefix, userID)

	err := rc.client.Set(ctx, key, permissionsJSON, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set user permissions in Redis: %w", err)
	}

	slog.Debug("User permissions cached", "userID", userID, "ttl", ttl, "dataSize", len(permissionsJSON))
	return nil
}

// DeleteUserPermissions remove permissões de usuário do Redis
func (rc *RedisCache) DeleteUserPermissions(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("%suser_permissions:%d", rc.keyPrefix, userID)

	err := rc.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete user permissions from Redis: %w", err)
	}

	slog.Debug("User permissions cache invalidated", "userID", userID)
	return nil
}
