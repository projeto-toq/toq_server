package permissionservice

import (
	"context"
	"strings"

	permissionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/permission_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListPermissions retorna permissões com filtros e paginação.
func (p *permissionServiceImpl) ListPermissions(ctx context.Context, input ListPermissionsInput) (ListPermissionsOutput, error) {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ListPermissionsOutput{}, utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 20
	}

	tx, txErr := p.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.permission.list.tx_start_failed", "error", txErr)
		return ListPermissionsOutput{}, utils.InternalError("")
	}
	defer func() {
		_ = p.globalService.RollbackTransaction(ctx, tx)
	}()

	filter := permissionrepository.PermissionListFilter{
		Page:     input.Page,
		Limit:    input.Limit,
		Name:     utils.NormalizeSearchPattern(strings.TrimSpace(input.Name)),
		Resource: utils.NormalizeSearchPattern(strings.TrimSpace(input.Resource)),
		Action:   utils.NormalizeSearchPattern(strings.TrimSpace(input.Action)),
		IsActive: input.IsActive,
	}

	result, listErr := p.permissionRepository.ListPermissions(ctx, tx, filter)
	if listErr != nil {
		utils.SetSpanError(ctx, listErr)
		logger.Error("permission.permission.list.repo_error", "error", listErr)
		return ListPermissionsOutput{}, utils.InternalError("")
	}

	output := ListPermissionsOutput{
		Permissions: result.Permissions,
		Total:       result.Total,
		Page:        filter.Page,
		Limit:       filter.Limit,
	}

	return output, nil
}
