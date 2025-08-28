package permissionservice

import (
	"context"
	"fmt"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// HasPermission verifica se o usuário tem uma permissão específica
func (p *permissionServiceImpl) HasPermission(ctx context.Context, userID int64, resource, action string, permContext *permissionmodel.PermissionContext) (bool, error) {
	if userID <= 0 {
		return false, fmt.Errorf("invalid user ID: %d", userID)
	}

	if resource == "" || action == "" {
		return false, fmt.Errorf("resource and action cannot be empty")
	}

	slog.Debug("Checking permission", "userID", userID, "resource", resource, "action", action)

	// Buscar permissões do usuário diretamente do banco
	userPermissions, err := p.getUserPermissionsFromDB(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user permissions: %w", err)
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
} // getUserPermissionsFromDB busca permissões do usuário no banco de dados
func (p *permissionServiceImpl) getUserPermissionsFromDB(ctx context.Context, userID int64) ([]permissionmodel.PermissionInterface, error) {
	// Usar o método do repositório para buscar todas as permissões do usuário
	permissions, err := p.permissionRepository.GetUserPermissions(ctx, nil, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions from repository: %w", err)
	}

	return permissions, nil
}
