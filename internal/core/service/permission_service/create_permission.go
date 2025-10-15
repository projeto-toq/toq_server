package permissionservice

import (
	"context"
	"strings"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreatePermission cria uma nova permissão no sistema.
func (p *permissionServiceImpl) CreatePermission(ctx context.Context, input CreatePermissionInput) (permissionmodel.PermissionInterface, error) {
	ctx, end, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	input.Name = strings.TrimSpace(input.Name)
	input.Resource = strings.TrimSpace(input.Resource)
	input.Action = strings.TrimSpace(input.Action)

	if input.Name == "" {
		return nil, utils.ValidationError("name", "permission name cannot be empty")
	}
	if input.Resource == "" {
		return nil, utils.ValidationError("resource", "permission resource cannot be empty")
	}
	if input.Action == "" {
		return nil, utils.ValidationError("action", "permission action cannot be empty")
	}

	// Verificar duplicidades em transação somente leitura
	roTx, txErr := p.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("permission.permission.create.tx_ro_start_failed", "error", txErr)
		return nil, utils.InternalError("")
	}

	existingByName, nameErr := p.permissionRepository.GetPermissionByName(ctx, roTx, input.Name)
	if nameErr != nil {
		utils.SetSpanError(ctx, nameErr)
		logger.Error("permission.permission.create.get_by_name_failed", "name", input.Name, "error", nameErr)
		_ = p.globalService.RollbackTransaction(ctx, roTx)
		return nil, utils.InternalError("")
	}
	if existingByName != nil {
		_ = p.globalService.RollbackTransaction(ctx, roTx)
		return nil, utils.ConflictError("permission name already exists")
	}

	existingByResource, resourceErr := p.permissionRepository.GetPermissionsByResourceAndAction(ctx, roTx, input.Resource, input.Action)
	if resourceErr != nil {
		utils.SetSpanError(ctx, resourceErr)
		logger.Error("permission.permission.create.get_by_resource_action_failed", "resource", input.Resource, "action", input.Action, "error", resourceErr)
		_ = p.globalService.RollbackTransaction(ctx, roTx)
		return nil, utils.InternalError("")
	}
	if len(existingByResource) > 0 {
		_ = p.globalService.RollbackTransaction(ctx, roTx)
		return nil, utils.ConflictError("permission resource/action already exists")
	}

	if commitErr := p.globalService.CommitTransaction(ctx, roTx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("permission.permission.create.tx_ro_commit_failed", "error", commitErr)
		return nil, utils.InternalError("")
	}

	// Criar domínio
	permission := permissionmodel.NewPermission()
	permission.SetName(input.Name)
	permission.SetResource(input.Resource)
	permission.SetAction(input.Action)
	permission.SetDescription(strings.TrimSpace(input.Description))
	if input.Conditions != nil {
		permission.SetConditions(input.Conditions)
	}
	permission.SetIsActive(true)

	// Transação de escrita
	tx, werr := p.globalService.StartTransaction(ctx)
	if werr != nil {
		utils.SetSpanError(ctx, werr)
		logger.Error("permission.permission.create.tx_start_failed", "resource", input.Resource, "action", input.Action, "error", werr)
		return nil, utils.InternalError("")
	}

	var opErr error
	defer func() {
		if opErr != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("permission.permission.create.tx_rollback_failed", "resource", input.Resource, "action", input.Action, "error", rbErr)
			}
		}
	}()

	if err = p.permissionRepository.CreatePermission(ctx, tx, permission); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("permission.permission.create.repo_error", "resource", input.Resource, "action", input.Action, "error", err)
		opErr = utils.InternalError("")
		return nil, opErr
	}

	if commitErr := p.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("permission.permission.create.tx_commit_failed", "permission_id", permission.GetID(), "error", commitErr)
		return nil, utils.InternalError("")
	}

	logger.Info("permission.permission.created", "permission_id", permission.GetID(), "resource", input.Resource, "action", input.Action)
	return permission, nil
}
