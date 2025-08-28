package permissionservice

import (
	"context"
	"fmt"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
)

// HasHTTPPermission verifica se o usuário tem permissão para um endpoint HTTP específico
func (p *permissionServiceImpl) HasHTTPPermission(ctx context.Context, userID int64, method, path string) (bool, error) {
	if userID <= 0 {
		return false, fmt.Errorf("invalid user ID: %d", userID)
	}

	if method == "" || path == "" {
		return false, fmt.Errorf("method and path cannot be empty")
	}

	slog.Debug("Checking HTTP permission", "userID", userID, "method", method, "path", path)

	// Mapear HTTP method+path para resource+action
	resource := "http"
	action := fmt.Sprintf("%s:%s", method, path)

	// Criar contexto básico para HTTP
	permContext := permissionmodel.NewPermissionContext(userID)
	permContext.AddMetadata("http_method", method)
	permContext.AddMetadata("http_path", path)

	// Usar o método HasPermission para verificar
	return p.HasPermission(ctx, userID, resource, action, permContext)
}
