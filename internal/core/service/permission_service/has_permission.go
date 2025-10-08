package permissionservice

import (
	"context"
	"encoding/json"
	"fmt"
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

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return false, utils.BadRequest("invalid user id")
	}

	if resource == "" || action == "" {
		return false, utils.BadRequest("invalid resource or action")
	}

	logger.Debug("permission.check.start", "user_id", userID, "resource", resource, "action", action)

	// Tentar buscar permissões do cache primeiro
	userPermissions, cacheHit, err := p.getUserPermissionsWithCache(ctx, userID)
	if err != nil {
		logger.Error("permission.check.permissions_load_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return false, utils.InternalError("")
	}

	// Verificar cada permissão
	evaluator := NewConditionEvaluator()
	evaluatePermissions := func(perms []permissionmodel.PermissionInterface) bool {
		for _, permission := range perms {
			if permission.GetResource() != resource || permission.GetAction() != action {
				continue
			}

			if permission.GetConditions() != nil {
				if permContext == nil {
					// Sem contexto mas com condições - negar
					continue
				}

				if evaluator.Evaluate(ctx, permission.GetConditions(), permContext) {
					logger.Info("permission.check.allowed", "user_id", userID, "resource", resource, "action", action, "permission_id", permission.GetID())
					return true
				}
				continue
			}

			// Permissão sem condições
			logger.Info("permission.check.allowed", "user_id", userID, "resource", resource, "action", action, "permission_id", permission.GetID())
			return true
		}

		return false
	}

	for attempt := 0; attempt < 2; attempt++ {
		if evaluatePermissions(userPermissions) {
			return true, nil
		}

		if attempt == 0 && cacheHit {
			logger.Info("permission.cache.refresh_on_miss", "user_id", userID, "resource", resource, "action", action)
			if err := p.RefreshUserPermissions(ctx, userID); err != nil {
				logger.Error("permission.cache.refresh_on_miss_failed", "user_id", userID, "resource", resource, "action", action, "error", err)
				utils.SetSpanError(ctx, err)
				return false, utils.InternalError("")
			}

			userPermissions, cacheHit, err = p.getUserPermissionsWithCache(ctx, userID)
			if err != nil {
				logger.Error("permission.cache.refresh_reload_failed", "user_id", userID, "resource", resource, "action", action, "error", err)
				utils.SetSpanError(ctx, err)
				return false, utils.InternalError("")
			}

			logger.Info("permission.cache.refresh_on_miss_completed", "user_id", userID, "resource", resource, "action", action, "permissions", len(userPermissions))
			continue
		}

		break
	}

	logger.Warn("permission.check.denied", "user_id", userID, "resource", resource, "action", action)
	return false, nil
}

// getUserPermissionsWithCache busca permissões com cache Redis
func (p *permissionServiceImpl) getUserPermissionsWithCache(ctx context.Context, userID int64) ([]permissionmodel.PermissionInterface, bool, error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	cacheKey := fmt.Sprintf("user_permissions:%d", userID)

	// Tentar buscar do cache Redis primeiro
	if p.cache != nil {
		cached, err := p.getUserPermissionsFromCache(ctx, cacheKey)
		if err == nil && cached != nil {
			logger.Debug("Permissions loaded from cache", "userID", userID, "count", len(cached))
			p.observeCacheOperation("user_permissions_lookup", "hit")
			return cached, true, nil
		}
		logger.Debug("Cache miss for user permissions", "userID", userID, "error", err)
		p.observeCacheOperation("user_permissions_lookup", "miss")
	} else {
		p.observeCacheOperation("user_permissions_lookup", "disabled")
	}

	// Cache miss ou erro - buscar do banco
	permissions, err := p.getUserPermissionsFromDB(ctx, userID)
	if err != nil {
		return nil, false, fmt.Errorf("db load error: %w", err)
	}

	// Armazenar no cache para próximas consultas
	if p.cache != nil {
		err = p.setUserPermissionsInCache(ctx, cacheKey, permissions)
		if err != nil {
			logger.Warn("Failed to cache user permissions", "userID", userID, "error", err)
			p.observeCacheOperation("user_permissions_store", "error")
			// Não falha - apenas loga o erro do cache
		} else {
			p.observeCacheOperation("user_permissions_store", "success")
		}
	} else {
		p.observeCacheOperation("user_permissions_store", "disabled")
	}

	logger.Debug("Permissions loaded from database", "userID", userID, "count", len(permissions))
	return permissions, false, nil
}

// getUserPermissionsFromCache busca permissões do cache Redis
func (p *permissionServiceImpl) getUserPermissionsFromCache(ctx context.Context, cacheKey string) ([]permissionmodel.PermissionInterface, error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	// Extrair userID da cacheKey (formato: "user_permissions:%d")
	parts := strings.Split(cacheKey, ":")
	if len(parts) != 2 {
		logger.Error("Invalid cache key format", "cache_key", cacheKey)
		return nil, fmt.Errorf("cache miss - invalid key format")
	}

	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		logger.Error("Failed to parse userID from cache key", "cache_key", cacheKey, "error", err)
		return nil, fmt.Errorf("cache miss - invalid user ID")
	}

	// Buscar do Redis usando os métodos existentes
	permissionsJSON, err := p.cache.GetUserPermissions(ctx, userID)
	if err != nil {
		// Cache miss é esperado e não é erro crítico
		logger.Debug("Cache miss for user permissions", "userID", userID, "error", err)
		return nil, fmt.Errorf("cache miss - %w", err)
	}

	// Deserializar JSON
	var cachedPermissions CachedUserPermissions
	if err := json.Unmarshal(permissionsJSON, &cachedPermissions); err != nil {
		logger.Error("Failed to unmarshal cached permissions", "userID", userID, "error", err)
		// Cache corrompido - limpar e forçar reload
		if deleteErr := p.cache.DeleteUserPermissions(ctx, userID); deleteErr != nil {
			logger.Warn("Failed to delete corrupted cache", "userID", userID, "error", deleteErr)
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

	logger.Debug("User permissions cache hit", "userID", userID, "count", len(permissions), "cached_at", cachedPermissions.CachedAt)
	return permissions, nil
}

// setUserPermissionsInCache armazena permissões no cache Redis
func (p *permissionServiceImpl) setUserPermissionsInCache(ctx context.Context, cacheKey string, permissions []permissionmodel.PermissionInterface) error {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	// Extrair userID da cacheKey (formato: "user_permissions:%d")
	parts := strings.Split(cacheKey, ":")
	if len(parts) != 2 {
		logger.Error("Invalid cache key format for storing", "cache_key", cacheKey)
		return fmt.Errorf("invalid cache key format")
	}

	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		logger.Error("Failed to parse userID from cache key for storing", "cache_key", cacheKey, "error", err)
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
		logger.Error("Failed to marshal permissions for cache", "userID", userID, "error", err)
		return fmt.Errorf("failed to serialize permissions: %w", err)
	}

	// Armazenar no Redis usando TTL padrão do cache
	// O TTL é definido internamente pelo RedisCache (15 minutos por padrão)
	err = p.cache.SetUserPermissions(ctx, userID, permissionsJSON, 15*time.Minute)
	if err != nil {
		logger.Error("Failed to store permissions in cache", "userID", userID, "error", err)
		return fmt.Errorf("failed to store in cache: %w", err)
	}

	logger.Debug("User permissions cached successfully", "userID", userID, "count", len(permissions), "dataSize", len(permissionsJSON))
	return nil
}

// getUserPermissionsFromDB busca permissões do usuário no banco de dados
func (p *permissionServiceImpl) getUserPermissionsFromDB(ctx context.Context, userID int64) ([]permissionmodel.PermissionInterface, error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.check.tx_start_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	// Rollback on error
	var retErr error
	defer func() {
		if retErr != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.check.tx_rollback_failed", "user_id", userID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	// Usar o método do repositório para buscar todas as permissões do usuário
	permissions, qerr := p.permissionRepository.GetUserPermissions(ctx, tx, userID)
	if qerr != nil {
		logger.Error("permission.check.db_failed", "user_id", userID, "error", qerr)
		utils.SetSpanError(ctx, qerr)
		retErr = utils.InternalError("")
		return nil, retErr
	}

	// Commit the transaction
	if cmErr := p.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		logger.Error("permission.check.tx_commit_failed", "user_id", userID, "error", cmErr)
		utils.SetSpanError(ctx, cmErr)
		retErr = utils.InternalError("")
		return nil, retErr
	}

	return permissions, nil
}

// ClearUserPermissionsCache remove as permissões do usuário do cache
func (p *permissionServiceImpl) ClearUserPermissionsCache(ctx context.Context, userID int64) error {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	if p.cache == nil {
		logger.Debug("permission.cache.clear.skip", "user_id", userID)
		p.observeCacheOperation("user_permissions_clear", "disabled")
		return nil
	}

	err := p.cache.DeleteUserPermissions(ctx, userID)
	if err != nil {
		logger.Error("permission.cache.clear_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		p.observeCacheOperation("user_permissions_clear", "error")
		// Retornar erro de infraestrutura de forma padronizada
		return utils.InternalError("")
	}

	logger.Info("permission.cache.invalidated", "user_id", userID)
	p.observeCacheOperation("user_permissions_clear", "success")
	return nil
}
