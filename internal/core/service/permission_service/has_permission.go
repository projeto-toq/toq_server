package permissionservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// CachedUserPermissions representa as permissões de usuário serializadas para cache
type CachedUserPermissions struct {
	CachedAt    time.Time                `json:"cached_at"`
	Permissions []PermissionSerializable `json:"permissions"`
}

// PermissionSerializable representa uma permissão serializada para JSON
type PermissionSerializable struct {
	ID          int64                  `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Resource    string                 `json:"resource"`
	Action      string                 `json:"action"`
	Conditions  map[string]interface{} `json:"conditions,omitempty"`
}

// HasPermission verifica se o usuário tem uma permissão específica
func (p *permissionServiceImpl) HasPermission(ctx context.Context, userID int64, resource, action string, permContext *permissionmodel.PermissionContext) (bool, error) {
	// Tracing da operação
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	if userID <= 0 {
		return false, utils.BadRequest("invalid user id")
	}

	if resource == "" || action == "" {
		return false, utils.BadRequest("invalid resource or action")
	}

	slog.Debug("permission.check.start", "user_id", userID, "resource", resource, "action", action)

	// Tentar buscar permissões do cache primeiro
	userPermissions, err := p.getUserPermissionsWithCache(ctx, userID)
	if err != nil {
		slog.Error("permission.check.permissions_load_failed", "user_id", userID, "error", err)
		return false, utils.InternalError("")
	}

	// Verificar cada permissão
	evaluator := NewConditionEvaluator()

	for _, permission := range userPermissions {
		if permission.GetResource() == resource && permission.GetAction() == action {
			// Verificar condições se existirem
			if permission.GetConditions() != nil {
				if permContext == nil {
					// Sem contexto mas com condições - negar
					continue
				}

				if evaluator.Evaluate(permission.GetConditions(), permContext) {
					slog.Info("permission.check.allowed", "user_id", userID, "resource", resource, "action", action, "permission_id", permission.GetID())
					return true, nil
				}
			} else {
				// Sem condições - permitir
				slog.Info("permission.check.allowed", "user_id", userID, "resource", resource, "action", action, "permission_id", permission.GetID())
				return true, nil
			}
		}
	}

	slog.Warn("permission.check.denied", "user_id", userID, "resource", resource, "action", action)
	return false, nil
}

// getUserPermissionsWithCache busca permissões com cache Redis
func (p *permissionServiceImpl) getUserPermissionsWithCache(ctx context.Context, userID int64) ([]permissionmodel.PermissionInterface, error) {
	cacheKey := fmt.Sprintf("user_permissions:%d", userID)

	// Tentar buscar do cache Redis primeiro
	if p.cache != nil {
		cached, err := p.getUserPermissionsFromCache(ctx, cacheKey)
		if err == nil && cached != nil {
			slog.Debug("Permissions loaded from cache", "userID", userID, "count", len(cached))
			return cached, nil
		}
		slog.Debug("Cache miss for user permissions", "userID", userID, "error", err)
	}

	// Cache miss ou erro - buscar do banco
	permissions, err := p.getUserPermissionsFromDB(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("db load error: %w", err)
	}

	// Armazenar no cache para próximas consultas
	if p.cache != nil {
		err = p.setUserPermissionsInCache(ctx, cacheKey, permissions)
		if err != nil {
			slog.Warn("Failed to cache user permissions", "userID", userID, "error", err)
			// Não falha - apenas loga o erro do cache
		}
	}

	slog.Debug("Permissions loaded from database", "userID", userID, "count", len(permissions))
	return permissions, nil
}

// getUserPermissionsFromCache busca permissões do cache Redis
func (p *permissionServiceImpl) getUserPermissionsFromCache(ctx context.Context, cacheKey string) ([]permissionmodel.PermissionInterface, error) {
	// Extrair userID da cacheKey (formato: "user_permissions:%d")
	parts := strings.Split(cacheKey, ":")
	if len(parts) != 2 {
		slog.Error("Invalid cache key format", "cache_key", cacheKey)
		return nil, fmt.Errorf("cache miss - invalid key format")
	}

	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		slog.Error("Failed to parse userID from cache key", "cache_key", cacheKey, "error", err)
		return nil, fmt.Errorf("cache miss - invalid user ID")
	}

	// Buscar do Redis usando os métodos existentes
	permissionsJSON, err := p.cache.GetUserPermissions(ctx, userID)
	if err != nil {
		// Cache miss é esperado e não é erro crítico
		slog.Debug("Cache miss for user permissions", "userID", userID, "error", err)
		return nil, fmt.Errorf("cache miss - %w", err)
	}

	// Deserializar JSON
	var cachedPermissions CachedUserPermissions
	if err := json.Unmarshal(permissionsJSON, &cachedPermissions); err != nil {
		slog.Error("Failed to unmarshal cached permissions", "userID", userID, "error", err)
		// Cache corrompido - limpar e forçar reload
		if deleteErr := p.cache.DeleteUserPermissions(ctx, userID); deleteErr != nil {
			slog.Warn("Failed to delete corrupted cache", "userID", userID, "error", deleteErr)
		}
		return nil, fmt.Errorf("cache miss - corrupted data")
	}

	// Converter structs serializáveis de volta para interfaces
	permissions := make([]permissionmodel.PermissionInterface, 0, len(cachedPermissions.Permissions))
	for _, perm := range cachedPermissions.Permissions {
		permission := permissionmodel.NewPermission()
		permission.SetID(perm.ID)
		permission.SetName(perm.Name)
		permission.SetDescription(perm.Description)
		permission.SetResource(perm.Resource)
		permission.SetAction(perm.Action)
		if perm.Conditions != nil {
			permission.SetConditions(perm.Conditions)
		}
		permissions = append(permissions, permission)
	}

	slog.Debug("User permissions cache hit", "userID", userID, "count", len(permissions), "cached_at", cachedPermissions.CachedAt)
	return permissions, nil
}

// setUserPermissionsInCache armazena permissões no cache Redis
func (p *permissionServiceImpl) setUserPermissionsInCache(ctx context.Context, cacheKey string, permissions []permissionmodel.PermissionInterface) error {
	// Extrair userID da cacheKey (formato: "user_permissions:%d")
	parts := strings.Split(cacheKey, ":")
	if len(parts) != 2 {
		slog.Error("Invalid cache key format for storing", "cache_key", cacheKey)
		return fmt.Errorf("invalid cache key format")
	}

	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		slog.Error("Failed to parse userID from cache key for storing", "cache_key", cacheKey, "error", err)
		return fmt.Errorf("invalid user ID in cache key: %w", err)
	}

	// Converter interfaces para structs serializáveis
	serializable := make([]PermissionSerializable, 0, len(permissions))
	for _, perm := range permissions {
		permStruct := PermissionSerializable{
			ID:          perm.GetID(),
			Name:        perm.GetName(),
			Description: perm.GetDescription(),
			Resource:    perm.GetResource(),
			Action:      perm.GetAction(),
			Conditions:  perm.GetConditions(),
		}
		serializable = append(serializable, permStruct)
	}

	// Criar estrutura com timestamp
	cachedPermissions := CachedUserPermissions{
		CachedAt:    time.Now(),
		Permissions: serializable,
	}

	// Serializar para JSON
	permissionsJSON, err := json.Marshal(cachedPermissions)
	if err != nil {
		slog.Error("Failed to marshal permissions for cache", "userID", userID, "error", err)
		return fmt.Errorf("failed to serialize permissions: %w", err)
	}

	// Armazenar no Redis usando TTL padrão do cache
	// O TTL é definido internamente pelo RedisCache (15 minutos por padrão)
	err = p.cache.SetUserPermissions(ctx, userID, permissionsJSON, 15*time.Minute)
	if err != nil {
		slog.Error("Failed to store permissions in cache", "userID", userID, "error", err)
		return fmt.Errorf("failed to store in cache: %w", err)
	}

	slog.Debug("User permissions cached successfully", "userID", userID, "count", len(permissions), "dataSize", len(permissionsJSON))
	return nil
}

// getUserPermissionsFromDB busca permissões do usuário no banco de dados
func (p *permissionServiceImpl) getUserPermissionsFromDB(ctx context.Context, userID int64) ([]permissionmodel.PermissionInterface, error) {
	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("permission.check.tx_start_failed", "user_id", userID, "error", err)
		return nil, utils.InternalError("")
	}
	// Rollback on error
	var retErr error
	defer func() {
		if retErr != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("permission.check.tx_rollback_failed", "user_id", userID, "error", rbErr)
			}
		}
	}()

	// Usar o método do repositório para buscar todas as permissões do usuário
	permissions, qerr := p.permissionRepository.GetUserPermissions(ctx, tx, userID)
	if qerr != nil {
		slog.Error("permission.check.db_failed", "user_id", userID, "error", qerr)
		retErr = utils.InternalError("")
		return nil, retErr
	}

	// Commit the transaction
	if cmErr := p.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		slog.Error("permission.check.tx_commit_failed", "user_id", userID, "error", cmErr)
		retErr = utils.InternalError("")
		return nil, retErr
	}

	return permissions, nil
}

// ClearUserPermissionsCache remove as permissões do usuário do cache
func (p *permissionServiceImpl) ClearUserPermissionsCache(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	err := p.cache.DeleteUserPermissions(ctx, userID)
	if err != nil {
		slog.Warn("permission.cache.clear_failed", "user_id", userID, "error", err)
		// Retornar erro de infraestrutura de forma padronizada
		return utils.InternalError("")
	}

	slog.Info("permission.cache.invalidated", "user_id", userID)
	return nil
}
