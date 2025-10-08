package permissionservice

import (
	"context"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetRolePermissions retorna todas as permissões de um role
func (p *permissionServiceImpl) GetRolePermissions(ctx context.Context, roleID int64) ([]permissionmodel.PermissionInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if roleID <= 0 {
		return nil, utils.BadRequest("invalid role id")
	}

	permissions, err := p.permissionRepository.GetGrantedPermissionsByRoleID(ctx, nil, roleID)
	if err != nil {
		logger.Error("permission.role.permissions.db_failed", "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return permissions, nil
}
