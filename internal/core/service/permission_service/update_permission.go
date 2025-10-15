package permissionservice

import (
	"context"
	"strings"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdatePermission atualiza campos permitidos de uma permissão existente.
func (p *permissionServiceImpl) UpdatePermission(ctx context.Context, input UpdatePermissionInput) (permissionmodel.PermissionInterface, error) {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if input.ID <= 0 {
		return nil, utils.BadRequest("invalid permission id")
	}

	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return nil, utils.ValidationError("name", "permission name cannot be empty")
	}

	tx, txErr := p.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.permission.update.tx_start_failed", "permission_id", input.ID, "error", txErr)
		return nil, utils.InternalError("")
	}

	var opErr error
	defer func() {
		if opErr != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("permission.permission.update.tx_rollback_failed", "permission_id", input.ID, "error", rbErr)
			}
		}
	}()

	existing, repoErr := p.permissionRepository.GetPermissionByID(ctx, tx, input.ID)
	if repoErr != nil {
		utils.SetSpanError(ctx, repoErr)
		logger.Error("permission.permission.update.repo_get_error", "permission_id", input.ID, "error", repoErr)
		opErr = utils.InternalError("")
		return nil, opErr
	}
	if existing == nil {
		opErr = utils.NotFoundError("permission")
		return nil, opErr
	}

	// Validar nome duplicado
	dup, dupErr := p.permissionRepository.GetPermissionByName(ctx, tx, input.Name)
	if dupErr != nil {
		utils.SetSpanError(ctx, dupErr)
		logger.Error("permission.permission.update.get_by_name_failed", "name", input.Name, "error", dupErr)
		opErr = utils.InternalError("")
		return nil, opErr
	}
	if dup != nil && dup.GetID() != existing.GetID() {
		opErr = utils.ConflictError("permission name already exists")
		return nil, opErr
	}

	// Capturar usuários impactados
	roleIDs, roleErr := p.permissionRepository.GetRoleIDsByPermissionID(ctx, tx, existing.GetID())
	if roleErr != nil {
		utils.SetSpanError(ctx, roleErr)
		logger.Error("permission.permission.update.get_role_ids_failed", "permission_id", input.ID, "error", roleErr)
		opErr = utils.InternalError("")
		return nil, opErr
	}

	userIDSet := make(map[int64]struct{})
	for _, roleID := range roleIDs {
		ids, usersErr := p.permissionRepository.GetActiveUserIDsByRoleID(ctx, tx, roleID)
		if usersErr != nil {
			utils.SetSpanError(ctx, usersErr)
			logger.Error("permission.permission.update.get_role_users_failed", "permission_id", input.ID, "role_id", roleID, "error", usersErr)
			opErr = utils.InternalError("")
			return nil, opErr
		}
		for _, uid := range ids {
			userIDSet[uid] = struct{}{}
		}
	}

	existing.SetName(input.Name)
	existing.SetDescription(strings.TrimSpace(input.Description))
	if input.IsActive != nil {
		existing.SetIsActive(*input.IsActive)
	}
	if input.Conditions != nil {
		existing.SetConditions(input.Conditions)
	}

	if updateErr := p.permissionRepository.UpdatePermission(ctx, tx, existing); updateErr != nil {
		utils.SetSpanError(ctx, updateErr)
		logger.Error("permission.permission.update.repo_error", "permission_id", input.ID, "error", updateErr)
		opErr = utils.InternalError("")
		return nil, opErr
	}

	if commitErr := p.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("permission.permission.update.tx_commit_failed", "permission_id", input.ID, "error", commitErr)
		return nil, utils.InternalError("")
	}

	for uid := range userIDSet {
		p.invalidateUserCacheSafe(ctx, uid, "update_permission")
	}

	logger.Info("permission.permission.updated", "permission_id", existing.GetID(), "active", existing.GetIsActive())
	return existing, nil
}
