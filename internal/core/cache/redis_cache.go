package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

// RedisCache implementa CacheInterface usando Redis
type RedisCache struct {
	client        *redis.Client
	globalService globalservice.GlobalServiceInterface
	defaultTTL    time.Duration
	keyPrefix     string
	metrics       metricsport.MetricsPortInterface
}

const (
	cacheOperationGet    = "get"
	cacheOperationSet    = "set"
	cacheOperationDelete = "delete"
	cacheOperationScan   = "scan"

	cacheResultHit     = "hit"
	cacheResultMiss    = "miss"
	cacheResultExpired = "expired"
	cacheResultError   = "error"
	cacheResultSuccess = "success"
)

// PermissionCacheEntry representa uma entrada de cache para o sistema de permissões moderno
type PermissionCacheEntry struct {
	Allowed   bool                     `json:"allowed"`
	Valid     bool                     `json:"valid"`
	CreatedAt time.Time                `json:"created_at"`
	UserID    int64                    `json:"user_id"`
	Action    string                   `json:"action"`
	RoleSlug  permissionmodel.RoleSlug `json:"role_slug"`
	ExpiresAt *time.Time               `json:"expires_at,omitempty"`
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
func NewRedisCache(redisURL string, globalService globalservice.GlobalServiceInterface, metrics metricsport.MetricsPortInterface) (CacheInterface, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	// Força Protocol 2 (RESP2) para evitar warning de maint_notifications
	// Redis 8.0.5 com RESP3 tenta usar maint_notifications que não está disponível
	opts.Protocol = 2

	client := redis.NewClient(opts)

	if err := redisotel.InstrumentTracing(client); err != nil {
		slog.Error("Failed to instrument Redis tracing", "error", err)
	}
	if err := redisotel.InstrumentMetrics(client); err != nil {
		slog.Error("Failed to instrument Redis metrics", "error", err)
	}

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
		metrics:       metrics,
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
func (rc *RedisCache) Get(ctx context.Context, userID int64, action string) (allowed bool, valid bool, err error) {
	key := rc.buildPermissionKey(userID, action)

	result, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Cache miss - não é erro
			rc.recordCacheOperation(cacheOperationGet, cacheResultMiss)
			return false, false, nil
		}
		rc.recordCacheOperation(cacheOperationGet, cacheResultError)
		slog.Error("Redis cache get error", "key", key, "error", err)
		return false, false, err
	}

	var entry PermissionCacheEntry
	if err := json.Unmarshal([]byte(result), &entry); err != nil {
		rc.recordCacheOperation(cacheOperationGet, cacheResultError)
		slog.Error("Failed to unmarshal cache entry", "key", key, "error", err)
		return false, false, err
	}

	// Verificar se o cache não expirou
	if time.Since(entry.CreatedAt) > rc.defaultTTL {
		// Cache expirado
		rc.recordCacheOperation(cacheOperationGet, cacheResultExpired)
		if err := rc.client.Del(ctx, key).Err(); err != nil {
			rc.recordCacheOperation(cacheOperationDelete, cacheResultError)
			slog.Error("Failed to delete expired cache entry", "key", key, "error", err)
		} else {
			rc.recordCacheOperation(cacheOperationDelete, cacheResultSuccess)
		}
		return false, false, nil
	}

	rc.recordCacheOperation(cacheOperationGet, cacheResultHit)
	slog.Debug("Cache hit", "key", key, "allowed", entry.Allowed)
	return entry.Allowed, entry.Valid, nil
}

// Set armazena uma entrada no cache para o sistema de permissões moderno
func (rc *RedisCache) Set(ctx context.Context, userID int64, action string, allowed bool, roleSlug permissionmodel.RoleSlug) error {
	key := rc.buildPermissionKey(userID, action)

	entry := PermissionCacheEntry{
		Allowed:   allowed,
		Valid:     true,
		CreatedAt: time.Now(),
		UserID:    userID,
		Action:    action,
		RoleSlug:  roleSlug,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		rc.recordCacheOperation(cacheOperationSet, cacheResultError)
		slog.Error("Failed to marshal cache entry", "key", key, "error", err)
		return err
	}

	err = rc.client.Set(ctx, key, data, rc.defaultTTL).Err()
	if err != nil {
		rc.recordCacheOperation(cacheOperationSet, cacheResultError)
		slog.Error("Failed to set cache entry", "key", key, "error", err)
		return err
	}

	rc.recordCacheOperation(cacheOperationSet, cacheResultSuccess)
	slog.Debug("Cache set", "key", key, "allowed", allowed)
	return nil
}

// GetRole recupera informações de role do cache
func (rc *RedisCache) GetRole(ctx context.Context, roleSlug permissionmodel.RoleSlug) (roleID int64, valid bool, err error) {
	key := rc.buildRoleKey(roleSlug)

	result, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			rc.recordCacheOperation(cacheOperationGet, cacheResultMiss)
			return 0, false, nil
		}
		rc.recordCacheOperation(cacheOperationGet, cacheResultError)
		slog.Error("Redis cache get role error", "key", key, "error", err)
		return 0, false, err
	}

	var entry RoleCacheEntry
	if err := json.Unmarshal([]byte(result), &entry); err != nil {
		rc.recordCacheOperation(cacheOperationGet, cacheResultError)
		slog.Error("Failed to unmarshal role cache entry", "key", key, "error", err)
		return 0, false, err
	}

	// Verificar se o cache não expirou
	if time.Since(entry.CachedAt) > rc.defaultTTL {
		rc.recordCacheOperation(cacheOperationGet, cacheResultExpired)
		if err := rc.client.Del(ctx, key).Err(); err != nil {
			rc.recordCacheOperation(cacheOperationDelete, cacheResultError)
			slog.Error("Failed to delete expired role cache entry", "key", key, "error", err)
		} else {
			rc.recordCacheOperation(cacheOperationDelete, cacheResultSuccess)
		}
		return 0, false, nil
	}

	rc.recordCacheOperation(cacheOperationGet, cacheResultHit)
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
		rc.recordCacheOperation(cacheOperationSet, cacheResultError)
		slog.Error("Failed to marshal role cache entry", "key", key, "error", err)
		return err
	}

	err = rc.client.Set(ctx, key, data, rc.defaultTTL).Err()
	if err != nil {
		rc.recordCacheOperation(cacheOperationSet, cacheResultError)
		slog.Error("Failed to set role cache entry", "key", key, "error", err)
		return err
	}

	rc.recordCacheOperation(cacheOperationSet, cacheResultSuccess)
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
		rc.recordCacheOperation(cacheOperationScan, cacheResultError)
		slog.Error("Error scanning cache keys", "error", err)
		return
	}

	if len(keys) > 0 {
		deleted, err := rc.client.Del(ctx, keys...).Result()
		if err != nil {
			rc.recordCacheOperation(cacheOperationDelete, cacheResultError)
			slog.Error("Error deleting cache keys", "error", err)
			return
		}
		rc.recordCacheOperation(cacheOperationDelete, cacheResultSuccess)
		slog.Info("Cache cleaned", "deleted_keys", deleted)
	} else {
		slog.Debug("No cache keys to clean")
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
		rc.recordCacheOperation(cacheOperationScan, cacheResultError)
		slog.Error("Error scanning cache keys by user", "user_id", userID, "error", err)
		return
	}

	if len(keys) > 0 {
		deleted, err := rc.client.Del(ctx, keys...).Result()
		if err != nil {
			rc.recordCacheOperation(cacheOperationDelete, cacheResultError)
			slog.Error("Error deleting cache keys by user", "user_id", userID, "error", err)
			return
		}
		rc.recordCacheOperation(cacheOperationDelete, cacheResultSuccess)
		slog.Debug("Cache cleaned by user", "user_id", userID, "deleted_keys", deleted)
	}
}

// buildPermissionKey constrói a chave do cache para permissões
func (rc *RedisCache) buildPermissionKey(userID int64, action string) string {
	return fmt.Sprintf("%spermission:user:%d:action:%s", rc.keyPrefix, userID, action)
}

// buildRoleKey constrói a chave do cache para roles
func (rc *RedisCache) buildRoleKey(roleSlug permissionmodel.RoleSlug) string {
	return fmt.Sprintf("%srole:slug:%s", rc.keyPrefix, roleSlug)
}

func (rc *RedisCache) recordCacheOperation(operation, result string) {
	if rc.metrics == nil {
		return
	}
	rc.metrics.IncrementCacheOperations(operation, result)
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
			rc.recordCacheOperation(cacheOperationGet, cacheResultMiss)
			return nil, fmt.Errorf("cache miss - user permissions not found for user %d", userID)
		}
		rc.recordCacheOperation(cacheOperationGet, cacheResultError)
		return nil, fmt.Errorf("failed to get user permissions from Redis: %w", err)
	}

	rc.recordCacheOperation(cacheOperationGet, cacheResultHit)
	slog.Debug("User permissions cache hit", "userID", userID, "dataSize", len(result))
	return []byte(result), nil
}

// SetUserPermissions armazena permissões de usuário no Redis
func (rc *RedisCache) SetUserPermissions(ctx context.Context, userID int64, permissionsJSON []byte, ttl time.Duration) error {
	key := fmt.Sprintf("%suser_permissions:%d", rc.keyPrefix, userID)

	err := rc.client.Set(ctx, key, permissionsJSON, ttl).Err()
	if err != nil {
		rc.recordCacheOperation(cacheOperationSet, cacheResultError)
		return fmt.Errorf("failed to set user permissions in Redis: %w", err)
	}

	rc.recordCacheOperation(cacheOperationSet, cacheResultSuccess)
	slog.Debug("User permissions cached", "userID", userID, "ttl", ttl, "dataSize", len(permissionsJSON))
	return nil
}

// DeleteUserPermissions remove permissões de usuário do Redis
func (rc *RedisCache) DeleteUserPermissions(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("%suser_permissions:%d", rc.keyPrefix, userID)

	err := rc.client.Del(ctx, key).Err()
	if err != nil {
		rc.recordCacheOperation(cacheOperationDelete, cacheResultError)
		return fmt.Errorf("failed to delete user permissions from Redis: %w", err)
	}

	rc.recordCacheOperation(cacheOperationDelete, cacheResultSuccess)
	slog.Debug("User permissions cache invalidated", "userID", userID)
	return nil
}
