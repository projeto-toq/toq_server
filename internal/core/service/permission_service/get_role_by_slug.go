package permissionservice

import (
	"context"
	"database/sql"

	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetRoleBySlug busca um role pelo slug (sem transação - uso direto)
func (p *permissionServiceImpl) GetRoleBySlug(ctx context.Context, slug permissionmodel.RoleSlug) (permissionmodel.RoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if slug == "" {
		return nil, utils.BadRequest("invalid role slug")
	}

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.role.get_by_slug.tx_start_failed", "slug", slug, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.role.get_by_slug.tx_rollback_failed", "slug", slug, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	role, derr := p.GetRoleBySlugWithTx(ctx, tx, slug)
	if derr != nil {
		return nil, derr
	}

	// Commit the transaction
	if err = p.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.role.get_by_slug.tx_commit_failed", "slug", slug, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return role, nil
}

// GetRoleBySlugWithTx busca um role pelo slug (com transação - uso em fluxos)
func (p *permissionServiceImpl) GetRoleBySlugWithTx(ctx context.Context, tx *sql.Tx, slug permissionmodel.RoleSlug) (permissionmodel.RoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	role, err := p.permissionRepository.GetRoleBySlug(ctx, tx, slug.String())
	if err != nil {
		logger.Error("permission.role.get_by_slug.db_failed", "slug", slug, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return role, nil
}
