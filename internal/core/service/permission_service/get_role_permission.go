package permissionservice

import (
	"context"
	"database/sql"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetRolePermissionByID retrieves a role-permission relation by its identifier.
func (p *permissionServiceImpl) GetRolePermissionByID(ctx context.Context, rolePermissionID int64) (permissionmodel.RolePermissionInterface, error) {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if rolePermissionID <= 0 {
		return nil, utils.ValidationError("id", "id must be greater than zero")
	}

	tx, txErr := p.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.role_permission.detail.tx_start_failed", "role_permission_id", rolePermissionID, "error", txErr)
		return nil, utils.InternalError("")
	}

	rolePermission, repoErr := p.permissionRepository.GetRolePermissionByID(ctx, tx, rolePermissionID)
	if repoErr != nil {
		if repoErr == sql.ErrNoRows {
			_ = p.globalService.RollbackTransaction(ctx, tx)
			return nil, utils.NotFoundError("role_permission")
		}
		utils.SetSpanError(ctx, repoErr)
		logger.Error("permission.role_permission.detail.repo_error", "role_permission_id", rolePermissionID, "error", repoErr)
		_ = p.globalService.RollbackTransaction(ctx, tx)
		return nil, utils.InternalError("")
	}

	if cmErr := p.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("permission.role_permission.detail.tx_commit_failed", "role_permission_id", rolePermissionID, "error", cmErr)
		_ = p.globalService.RollbackTransaction(ctx, tx)
		return nil, utils.InternalError("")
	}

	return rolePermission, nil
}
