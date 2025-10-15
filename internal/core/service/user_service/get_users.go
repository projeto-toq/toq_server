package userservices

import (
	"context"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	userrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/user_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListUsers retorna usuários com filtros e paginação para o painel admin.
func (us *userService) ListUsers(ctx context.Context, input ListUsersInput) (ListUsersOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ListUsersOutput{}, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 20
	}

	tx, txErr := us.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("admin.users.list.tx_start_failed", "error", txErr)
		return ListUsersOutput{}, utils.InternalError("")
	}
	defer func() {
		_ = us.globalService.RollbackTransaction(ctx, tx)
	}()

	var statusPtr *permissionmodel.UserRoleStatus
	if input.RoleStatus != nil {
		statusCopy := *input.RoleStatus
		statusPtr = &statusCopy
	}

	filter := userrepository.ListUsersFilter{
		Page:             input.Page,
		Limit:            input.Limit,
		RoleName:         utils.NormalizeSearchPattern(input.RoleName),
		RoleSlug:         utils.NormalizeSearchPattern(input.RoleSlug),
		RoleStatus:       statusPtr,
		IsSystemRole:     input.IsSystemRole,
		FullName:         utils.NormalizeSearchPattern(input.FullName),
		CPF:              utils.NormalizeSearchPattern(input.CPF),
		Email:            utils.NormalizeSearchPattern(input.Email),
		PhoneNumber:      utils.NormalizeSearchPattern(input.PhoneNumber),
		Deleted:          input.Deleted,
		IDFrom:           input.IDFrom,
		IDTo:             input.IDTo,
		BornAtFrom:       input.BornAtFrom,
		BornAtTo:         input.BornAtTo,
		LastActivityFrom: input.LastActivityFrom,
		LastActivityTo:   input.LastActivityTo,
	}

	result, listErr := us.repo.ListUsersWithFilters(ctx, tx, filter)
	if listErr != nil {
		utils.SetSpanError(ctx, listErr)
		logger.Error("admin.users.list.repo_error", "error", listErr)
		return ListUsersOutput{}, utils.InternalError("")
	}

	output := ListUsersOutput{
		Users: result.Users,
		Total: result.Total,
		Page:  filter.Page,
		Limit: filter.Limit,
	}

	return output, nil
}
