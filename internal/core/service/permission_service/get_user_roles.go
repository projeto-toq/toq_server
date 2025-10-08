package permissionservice

import (
	"context"
	"database/sql"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetUserRoles returns all roles of a user, independent of is_active
func (p *permissionServiceImpl) GetUserRoles(ctx context.Context, userID int64) ([]permissionmodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user_roles.tx_start_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user_roles.tx_rollback_failed", "user_id", userID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	// Busca todas as roles do usuário (ativas e inativas); a regra de negócio prevê apenas uma ativa.
	userRoles, err := p.permissionRepository.GetUserRolesByUserID(ctx, tx, userID)
	if err != nil {
		logger.Error("permission.user_roles.db_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	// Commit the transaction
	if err = p.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.user_roles.tx_commit_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return userRoles, nil
}

// GetUserRolesWithTx returns all roles of a user within a provided transaction (used in flows)
func (p *permissionServiceImpl) GetUserRolesWithTx(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	// Busca todas as roles do usuário (ativas e inativas); a regra de negócio prevê apenas uma ativa.
	userRoles, err := p.permissionRepository.GetUserRolesByUserID(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("permission.user_roles.db_failed", "user_id", userID, "error", err)
		return nil, utils.InternalError("")
	}

	return userRoles, nil
}
