package permissionservice

import (
	"context"
	"database/sql"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetUserRoles retorna todos os roles ativos de um usuário
func (p *permissionServiceImpl) GetUserRoles(ctx context.Context, userID int64) ([]permissionmodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("permission.user_roles.tx_start_failed", "user_id", userID, "error", err)
		return nil, utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := p.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("permission.user_roles.tx_rollback_failed", "user_id", userID, "error", rbErr)
			}
		}
	}()

	userRoles, err := p.permissionRepository.GetActiveUserRolesByUserID(ctx, tx, userID)
	if err != nil {
		slog.Error("permission.user_roles.db_failed", "user_id", userID, "error", err)
		return nil, utils.InternalError("")
	}

	// Commit the transaction
	if err = p.globalService.CommitTransaction(ctx, tx); err != nil {
		slog.Error("permission.user_roles.tx_commit_failed", "user_id", userID, "error", err)
		return nil, utils.InternalError("")
	}

	return userRoles, nil
}

// GetUserRolesWithTx retorna todos os roles ativos de um usuário (com transação - uso em fluxos)
func (p *permissionServiceImpl) GetUserRolesWithTx(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.UserRoleInterface, error) {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	userRoles, err := p.permissionRepository.GetActiveUserRolesByUserID(ctx, tx, userID)
	if err != nil {
		return nil, utils.InternalError("")
	}

	return userRoles, nil
}
