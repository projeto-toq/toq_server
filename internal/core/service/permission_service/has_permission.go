package permissionservice

import (
	"context"
	"fmt"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// HasPermission verifica se o usuário tem uma permissão específica
func (p *permissionServiceImpl) HasPermission(ctx context.Context, userID int64, resource, action string, permContext *permissionmodel.PermissionContext) (bool, error) {
	if userID <= 0 {
		return false, utils.ErrBadRequest
	}

	if resource == "" || action == "" {
		return false, utils.ErrBadRequest
	}

	slog.Debug("Checking permission", "userID", userID, "resource", resource, "action", action)

	// Tentar buscar permissões do cache primeiro
	userPermissions, err := p.getUserPermissionsWithCache(ctx, userID)
	if err != nil {
		return false, utils.ErrInternalServer
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
					slog.Debug("Permission granted with conditions", "permission_id", permission.GetID())
					return true, nil
				}
			} else {
				// Sem condições - permitir
				slog.Debug("Permission granted without conditions", "permission_id", permission.GetID())
				return true, nil
			}
		}
	}

	slog.Debug("Permission denied", "userID", userID, "resource", resource, "action", action)
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
		return nil, utils.ErrInternalServer
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
func (p *permissionServiceImpl) getUserPermissionsFromCache(_ context.Context, cacheKey string) ([]permissionmodel.PermissionInterface, error) {
	// TODO: Implementar cache otimizado para o novo sistema
	// Por enquanto, sempre retorna cache miss para forçar busca no banco
	slog.Debug("Cache temporarily disabled for user permissions", "cache_key", cacheKey)
	return nil, fmt.Errorf("cache miss - using database")
}

// setUserPermissionsInCache armazena permissões no cache Redis
func (p *permissionServiceImpl) setUserPermissionsInCache(_ context.Context, cacheKey string, permissions []permissionmodel.PermissionInterface) error {
	// TODO: Implementar cache otimizado para o novo sistema
	// Por enquanto, não faz cache para simplificar a migração
	slog.Debug("Cache temporarily disabled for storing user permissions", "cache_key", cacheKey, "count", len(permissions))
	return nil
}

// getUserPermissionsFromDB busca permissões do usuário no banco de dados
func (p *permissionServiceImpl) getUserPermissionsFromDB(ctx context.Context, userID int64) ([]permissionmodel.PermissionInterface, error) {
	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		return nil, utils.ErrInternalServer
	}

	// Usar o método do repositório para buscar todas as permissões do usuário
	permissions, err := p.permissionRepository.GetUserPermissions(ctx, tx, userID)
	if err != nil {
		p.globalService.RollbackTransaction(ctx, tx)
		return nil, utils.ErrInternalServer
	}

	// Commit the transaction
	err = p.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		p.globalService.RollbackTransaction(ctx, tx)
		return nil, utils.ErrInternalServer
	}

	return permissions, nil
}
