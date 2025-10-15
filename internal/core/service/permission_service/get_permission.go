package permissionservice

import (
	"context"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetPermissionByID recupera permissão pelo identificador.
func (p *permissionServiceImpl) GetPermissionByID(ctx context.Context, permissionID int64) (permissionmodel.PermissionInterface, error) {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if permissionID <= 0 {
		return nil, utils.BadRequest("invalid permission id")
	}

	tx, txErr := p.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.permission.get.tx_start_failed", "permission_id", permissionID, "error", txErr)
		return nil, utils.InternalError("")
	}
	defer func() {
		_ = p.globalService.RollbackTransaction(ctx, tx)
	}()

	permission, repoErr := p.permissionRepository.GetPermissionByID(ctx, tx, permissionID)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("permission.permission.get.repo_error", "permission_id", permissionID, "error", repoErr)
		return nil, utils.InternalError("")
	}
	if permission == nil {
		return nil, utils.NotFoundError("permission")
	}

	return permission, nil
}
