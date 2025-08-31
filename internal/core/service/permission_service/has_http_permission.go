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

	// Buscar informações do usuário para obter UserRoleID e RoleStatus
	userRole, err := p.permissionRepository.GetActiveUserRoleByUserID(ctx, nil, userID)
	if err != nil {
		slog.Error("Failed to get user active role", "userID", userID, "error", err)
		return false, err
	}

	if userRole == nil {
		slog.Warn("User has no active role", "userID", userID)
		return false, nil
	}

	// Mapear HTTP method+path para resource+action
	resource := "http"
	action := fmt.Sprintf("%s:%s", method, path)

	// Criar contexto com as informações completas do usuário
	permContext := permissionmodel.NewPermissionContext(userID, userRole.GetID(), userRole.GetStatus())
	permContext.AddMetadata("http_method", method)
	permContext.AddMetadata("http_path", path)

	// Usar o método HasPermission para verificar
	return p.HasPermission(ctx, userID, resource, action, permContext)
}
