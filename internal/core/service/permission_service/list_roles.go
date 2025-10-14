package permissionservice

import (
	"context"
	"strings"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	permissionrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/permission_repository"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListRolesInput define filtros e paginação para listagem de roles.
type ListRolesInput struct {
	Page         int
	Limit        int
	Name         string
	Slug         string
	IsSystemRole *bool
	IsActive     *bool
}

// ListRolesOutput agrega roles e metadados de paginação.
type ListRolesOutput struct {
	Roles []permissionmodel.RoleInterface
	Total int64
	Page  int
	Limit int
}

// UpdateRoleInput define atributos atualizáveis de um role.
type UpdateRoleInput struct {
	ID          int64
	Name        string
	Description string
	IsActive    bool
}

// ListRoles aplica filtros sobre roles disponíveis.
func (p *permissionServiceImpl) ListRoles(ctx context.Context, input ListRolesInput) (ListRolesOutput, error) {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return ListRolesOutput{}, utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := p.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.role.list.tx_start_failed", "error", txErr)
		return ListRolesOutput{}, utils.InternalError("")
	}
	defer func() {
		_ = p.globalService.RollbackTransaction(ctx, tx)
	}()

	repoFilter := permissionrepository.RoleListFilter{
		Page:         input.Page,
		Limit:        input.Limit,
		Name:         input.Name,
		Slug:         input.Slug,
		IsSystemRole: input.IsSystemRole,
		IsActive:     input.IsActive,
	}

	res, listErr := p.permissionRepository.ListRoles(ctx, tx, repoFilter)
	if listErr != nil {
		utils.SetSpanError(ctx, listErr)
		logger.Error("permission.role.list.repo_error", "error", listErr)
		return ListRolesOutput{}, utils.InternalError("")
	}

	output := ListRolesOutput{
		Roles: res.Roles,
		Total: res.Total,
		Page:  repoFilter.Page,
		Limit: repoFilter.Limit,
	}

	return output, nil
}

// GetRoleByID retorna role pelo identificador.
func (p *permissionServiceImpl) GetRoleByID(ctx context.Context, roleID int64) (permissionmodel.RoleInterface, error) {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if roleID <= 0 {
		return nil, utils.BadRequest("invalid role id")
	}

	tx, txErr := p.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.role.get.tx_start_failed", "role_id", roleID, "error", txErr)
		return nil, utils.InternalError("")
	}
	defer func() {
		_ = p.globalService.RollbackTransaction(ctx, tx)
	}()

	role, repoErr := p.permissionRepository.GetRoleByID(ctx, tx, roleID)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("permission.role.get.repo_error", "role_id", roleID, "error", repoErr)
		return nil, utils.InternalError("")
	}
	if role == nil {
		return nil, utils.NotFoundError("role")
	}

	return role, nil
}

// UpdateRole aplica alterações permitidas em um role existente.
func (p *permissionServiceImpl) UpdateRole(ctx context.Context, input UpdateRoleInput) (permissionmodel.RoleInterface, error) {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ID <= 0 {
		return nil, utils.BadRequest("invalid role id")
	}
	if strings.TrimSpace(input.Name) == "" {
		return nil, utils.ValidationError("name", "role name cannot be empty")
	}

	tx, txErr := p.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.role.update.tx_start_failed", "role_id", input.ID, "error", txErr)
		return nil, utils.InternalError("")
	}

	var opErr error
	defer func() {
		if opErr != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("permission.role.update.tx_rollback_failed", "role_id", input.ID, "error", rbErr)
			}
		}
	}()

	existing, repoErr := p.permissionRepository.GetRoleByID(ctx, tx, input.ID)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("permission.role.update.repo_get_error", "role_id", input.ID, "error", repoErr)
		opErr = utils.InternalError("")
		return nil, opErr
	}
	if existing == nil {
		opErr = utils.NotFoundError("role")
		return nil, opErr
	}

	originalSlug := existing.GetSlug()

	existing.SetName(input.Name)
	existing.SetDescription(input.Description)
	existing.SetIsActive(input.IsActive)

	if updateErr := p.permissionRepository.UpdateRole(ctx, tx, existing); updateErr != nil {
		utils.SetSpanError(ctx, updateErr)
		logger.Error("permission.role.update.repo_error", "role_id", input.ID, "error", updateErr)
		opErr = utils.InternalError("")
		return nil, opErr
	}

	if commitErr := p.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("permission.role.update.tx_commit_failed", "role_id", input.ID, "error", commitErr)
		return nil, utils.InternalError("")
	}

	logger.Info("permission.role.updated", "role_id", input.ID, "slug", originalSlug)
	return existing, nil
}

// DeleteRole desativa o role garantindo que não haja usuários ativos associados.
func (p *permissionServiceImpl) DeleteRole(ctx context.Context, roleID int64) error {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if roleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	tx, txErr := p.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.role.delete.tx_start_failed", "role_id", roleID, "error", txErr)
		return utils.InternalError("")
	}

	var opErr error
	defer func() {
		if opErr != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("permission.role.delete.tx_rollback_failed", "role_id", roleID, "error", rbErr)
			}
		}
	}()

	existing, repoErr := p.permissionRepository.GetRoleByID(ctx, tx, roleID)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("permission.role.delete.repo_get_error", "role_id", roleID, "error", repoErr)
		opErr = utils.InternalError("")
		return opErr
	}
	if existing == nil {
		opErr = utils.NotFoundError("role")
		return opErr
	}

	slug := existing.GetSlug()
	if slug == "admin" {
		opErr = derrors.ErrAdminRoleProtected
		return opErr
	}

	activeUsers, usersErr := p.permissionRepository.GetActiveUserIDsByRoleID(ctx, tx, roleID)
	if usersErr != nil {
		utils.SetSpanError(ctx, usersErr)
		logger.Error("permission.role.delete.repo_users_error", "role_id", roleID, "error", usersErr)
		opErr = utils.InternalError("")
		return opErr
	}
	if len(activeUsers) > 0 {
		opErr = derrors.ErrRoleDeletionHasUsers
		return opErr
	}

	existing.SetIsActive(false)
	if updateErr := p.permissionRepository.UpdateRole(ctx, tx, existing); updateErr != nil {
		utils.SetSpanError(ctx, updateErr)
		logger.Error("permission.role.delete.repo_update_error", "role_id", roleID, "error", updateErr)
		opErr = utils.InternalError("")
		return opErr
	}

	if commitErr := p.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("permission.role.delete.tx_commit_failed", "role_id", roleID, "error", commitErr)
		return utils.InternalError("")
	}

	logger.Info("permission.role.deactivated", "role_id", roleID, "slug", slug)
	return nil
}
