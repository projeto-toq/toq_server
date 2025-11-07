package permissionservice

import (
	"context"
	"database/sql"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetRoleByID retorna um papel pelo seu ID
func (p *permissionServiceImpl) GetRoleByID(ctx context.Context, roleID int64) (permissionmodel.RoleInterface, error) {
	// Tracing da operação
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if roleID <= 0 {
		return nil, utils.BadRequest("invalid role id")
	}

	logger.Debug("permission.role.fetch.start", "role_id", roleID)

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.role.fetch.tx_start_failed", "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	role, err := p.GetRoleByIDWithTx(ctx, tx, roleID)
	if err != nil {
		logger.Error("permission.role.fetch.db_failed", "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	// Commit the transaction
	if err = p.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.role.fetch.tx_commit_failed", "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	logger.Info("permission.role.fetch.fetched", "role_id", roleID)
	return role, nil
}

// GetRoleByIDWithTx retorna um papel pelo seu ID (com transação - uso em fluxos)
func (p *permissionServiceImpl) GetRoleByIDWithTx(ctx context.Context, tx *sql.Tx, roleID int64) (permissionmodel.RoleInterface, error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	role, err := p.permissionRepository.GetRoleByID(ctx, tx, roleID)
	if err != nil {
		logger.Error("permission.role.fetch.db_failed", "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return role, nil
}
