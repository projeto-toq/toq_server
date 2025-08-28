package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	// TODO: Replace with HTTP-based permission system
	// "github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/redis/go-redis/v9"
	// "google.golang.org/grpc"
)

// RedisCache implementa CacheInterface usando Redis
type RedisCache struct {
	client        *redis.Client
	globalService globalservice.GlobalServiceInterface
	defaultTTL    time.Duration
	keyPrefix     string
}

// CacheEntry representa uma entrada no cache
type CacheEntry struct {
	Allowed   bool               `json:"allowed"`
	Valid     bool               `json:"valid"`
	CreatedAt time.Time          `json:"created_at"`
	Method    string             `json:"method"`
	Role      usermodel.UserRole `json:"role"`
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

// Get recupera uma entrada do cache
func (rc *RedisCache) Get(ctx context.Context, fullMethod string, role usermodel.UserRole) (allowed bool, valid bool, err error) {
	key := rc.buildKey(fullMethod, role)

	slog.Debug("Cache lookup", "key", key, "method", fullMethod, "role", role)

	// Tentar buscar no cache
	result, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Cache miss - buscar dados e armazenar
			slog.Debug("Cache miss", "key", key)
			return rc.fetchAndCache(ctx, fullMethod, role, key)
		}
		slog.Error("Redis get error", "key", key, "error", err)
		// Em caso de erro do Redis, buscar diretamente (fallback)
		return rc.fetchFromService(ctx, fullMethod, role)
	}

	// Cache hit - deserializar
	var entry CacheEntry
	err = json.Unmarshal([]byte(result), &entry)
	if err != nil {
		slog.Error("Failed to unmarshal cache entry", "key", key, "error", err)
		// Em caso de erro de deserialização, buscar diretamente
		return rc.fetchFromService(ctx, fullMethod, role)
	}

	slog.Debug("Cache hit", "key", key, "allowed", entry.Allowed, "valid", entry.Valid)
	return entry.Allowed, entry.Valid, nil
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

// CleanByMethod limpa cache específico de um método
func (rc *RedisCache) CleanByMethod(ctx context.Context, fullMethod string) {
	pattern := rc.keyPrefix + fmt.Sprintf("method:%s:*", fullMethod)

	slog.Debug("Cleaning cache by method", "method", fullMethod, "pattern", pattern)

	iter := rc.client.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		slog.Error("Error scanning cache keys by method", "method", fullMethod, "error", err)
		return
	}

	if len(keys) > 0 {
		deleted, err := rc.client.Del(ctx, keys...).Result()
		if err != nil {
			slog.Error("Error deleting cache keys by method", "method", fullMethod, "error", err)
			return
		}
		slog.Debug("Cache cleaned by method", "method", fullMethod, "deleted_keys", deleted)
	}
}

// buildKey constrói a chave do cache
func (rc *RedisCache) buildKey(fullMethod string, role usermodel.UserRole) string {
	return fmt.Sprintf("%smethod:%s:role:%d", rc.keyPrefix, fullMethod, role)
}

// fetchAndCache busca dados do serviço e armazena no cache
func (rc *RedisCache) fetchAndCache(ctx context.Context, fullMethod string, role usermodel.UserRole, key string) (bool, bool, error) {
	allowed, valid, err := rc.fetchFromService(ctx, fullMethod, role)
	if err != nil {
		return false, false, err
	}

	// Armazenar no cache
	entry := CacheEntry{
		Allowed:   allowed,
		Valid:     valid,
		CreatedAt: time.Now(),
		Method:    fullMethod,
		Role:      role,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		slog.Error("Failed to marshal cache entry", "key", key, "error", err)
		// Mesmo com erro de cache, retornar os dados
		return allowed, valid, nil
	}

	err = rc.client.Set(ctx, key, data, rc.defaultTTL).Err()
	if err != nil {
		slog.Error("Failed to set cache entry", "key", key, "error", err)
		// Mesmo com erro de cache, retornar os dados
	} else {
		slog.Debug("Cache entry stored", "key", key, "ttl", rc.defaultTTL)
	}

	return allowed, valid, nil
}

// decodeFullMethod decodifica o método completo para service e method ID
// TODO: Implement HTTP-based method decoding
func (rc *RedisCache) decodeFullMethod(fullMethod string) (service usermodel.GRPCService, method uint8, err error) {
	// Temporary implementation for HTTP migration
	slog.Debug("decodeFullMethod temporarily disabled", "method", fullMethod)
	return usermodel.ServiceUserService, 0, nil
}

// getMethodId obtém o ID do método baseado no nome
// TODO: Implement HTTP-based method ID resolution
// func (rc *RedisCache) getMethodId(methods interface{}, name string) uint8 {
// 	// Temporary implementation for HTTP migration
// 	slog.Debug("getMethodId temporarily disabled", "name", name)
// 	return 0
// }

// fetchFromService busca dados diretamente do serviço
func (rc *RedisCache) fetchFromService(ctx context.Context, fullMethod string, role usermodel.UserRole) (bool, bool, error) {
	service, method, err := rc.decodeFullMethod(fullMethod)
	if err != nil {
		slog.Error("Error decoding full method", "method", fullMethod, "error", err)
		return false, false, err
	}

	slog.Debug("Fetching from service", "service", service, "method", method, "role", role)

	privilege, err := rc.globalService.GetPrivilegeForCache(ctx, service, method, role)
	if err != nil {
		slog.Error("Error getting privilege from service", "service", service, "method", method, "role", role, "error", err)
		return false, false, err
	}

	if privilege == nil {
		slog.Debug("No privilege found", "service", service, "method", method, "role", role)
		return false, true, nil
	}

	allowed := privilege.Allowed()
	slog.Debug("Privilege fetched from service", "service", service, "method", method, "role", role, "allowed", allowed)

	return allowed, true, nil
}

// Close fecha a conexão com Redis
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
