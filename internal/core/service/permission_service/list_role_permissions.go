package permissionservice

import (
	"context"

	permissionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/permission_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListRolePermissions lista associações role-permission com filtros fornecidos.
func (p *permissionServiceImpl) ListRolePermissions(ctx context.Context, input ListRolePermissionsInput) (ListRolePermissionsOutput, error) {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ListRolePermissionsOutput{}, utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.RoleID == nil && input.PermissionID == nil {
		return ListRolePermissionsOutput{}, utils.ValidationError("filters", "roleId or permissionId required")
	}

	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 20
	}

	tx, txErr := p.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.role_permission.list.tx_start_failed", "error", txErr)
		return ListRolePermissionsOutput{}, utils.InternalError("")
	}
	defer func() {
		_ = p.globalService.RollbackTransaction(ctx, tx)
	}()

	filter := permissionrepository.RolePermissionListFilter{
		Page:         input.Page,
		Limit:        input.Limit,
		RoleID:       input.RoleID,
		PermissionID: input.PermissionID,
		Granted:      input.Granted,
	}

	result, listErr := p.permissionRepository.ListRolePermissions(ctx, tx, filter)
	if listErr != nil {
		utils.SetSpanError(ctx, listErr)
		logger.Error("permission.role_permission.list.repo_error", "error", listErr)
		return ListRolePermissionsOutput{}, utils.InternalError("")
	}

	output := ListRolePermissionsOutput{
		RolePermissions: result.RolePermissions,
		Total:           result.Total,
		Page:            filter.Page,
		Limit:           filter.Limit,
	}

	return output, nil
}
