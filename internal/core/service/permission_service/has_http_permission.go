package permissionservice

import (
	"context"
	"fmt"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// HasHTTPPermission verifica se o usuário tem permissão para um endpoint HTTP específico
func (p *permissionServiceImpl) HasHTTPPermission(ctx context.Context, userID int64, method, path string) (bool, error) {
	if userID <= 0 {
		return false, utils.ErrBadRequest
	}

	if method == "" || path == "" {
		return false, utils.ErrBadRequest
	}

	slog.Debug("Checking HTTP permission", "userID", userID, "method", method, "path", path)

	// Buscar informações do usuário para obter UserRoleID e RoleStatus usando transação
	tx, txErr := p.globalService.StartTransaction(ctx)
	if txErr != nil {
		slog.Error("Failed to start transaction for HasHTTPPermission", "userID", userID, "error", txErr)
		return false, utils.ErrInternalServer
	}

	userRole, err := p.permissionRepository.GetActiveUserRoleByUserID(ctx, tx, userID)
	if err != nil {
		slog.Error("Failed to get user active role", "userID", userID, "error", err)
		_ = p.globalService.RollbackTransaction(ctx, tx)
		return false, err
	}

	if userRole == nil {
		slog.Warn("User has no active role", "userID", userID)
		if cerr := p.globalService.CommitTransaction(ctx, tx); cerr != nil {
			_ = p.globalService.RollbackTransaction(ctx, tx)
			return false, utils.ErrInternalServer
		}
		return false, nil
	}

	// Mapear HTTP method+path para resource+action
	resource := "http"
	action := fmt.Sprintf("%s:%s", method, path)

	// Criar contexto com as informações completas do usuário
	permContext := permissionmodel.NewPermissionContext(userID, userRole.GetID(), userRole.GetStatus())
	permContext.AddMetadata("http_method", method)
	permContext.AddMetadata("http_path", path)

	// Usar o método HasPermission para verificar (fora da consulta de role)
	if cerr := p.globalService.CommitTransaction(ctx, tx); cerr != nil {
		_ = p.globalService.RollbackTransaction(ctx, tx)
		return false, utils.ErrInternalServer
	}
	return p.HasPermission(ctx, userID, resource, action, permContext)
}
