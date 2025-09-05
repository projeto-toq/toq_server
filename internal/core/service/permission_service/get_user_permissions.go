package permissionservice

import (
	"context"
	"database/sql"
	"log/slog"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// GetUserPermissions retorna todas as permissões de um usuário
func (p *permissionServiceImpl) GetUserPermissions(ctx context.Context, userID int64) ([]permissionmodel.PermissionInterface, error) {
	// Tracing da operação
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	slog.Debug("permission.user.permissions.fetch.start", "user_id", userID)

	// Start transaction
	tx, err := p.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("permission.user.permissions.tx_start_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	permissions, err := p.permissionRepository.GetUserPermissions(ctx, tx, userID)
	if err != nil {
		slog.Error("permission.user.permissions.db_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	// Commit the transaction
	if err = p.globalService.CommitTransaction(ctx, tx); err != nil {
		slog.Error("permission.user.permissions.tx_commit_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	slog.Info("permission.user.permissions.fetched", "user_id", userID, "count", len(permissions))
	return permissions, nil
}

// GetUserPermissionsWithTx retorna todas as permissões de um usuário (com transação - uso em fluxos)
func (p *permissionServiceImpl) GetUserPermissionsWithTx(ctx context.Context, tx *sql.Tx, userID int64) ([]permissionmodel.PermissionInterface, error) {
	if userID <= 0 {
		return nil, utils.BadRequest("invalid user id")
	}

	permissions, err := p.permissionRepository.GetUserPermissions(ctx, tx, userID)
	if err != nil {
		slog.Error("permission.user.permissions.db_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("")
	}

	return permissions, nil
}
