package permissionservice

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// RemoveRoleFromUser remove um role de um usuário
func (p *permissionServiceImpl) RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error {
	// Tracing da operação
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	if roleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	slog.Debug("permission.role.remove.start", "user_id", userID, "role_id", roleID)

	// Buscar o UserRole existente
	userRole, err := p.permissionRepository.GetUserRoleByUserIDAndRoleID(ctx, nil, userID, roleID)
	if err != nil {
		slog.Error("permission.role.get_user_role_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	if userRole == nil {
		return utils.NotFoundError("UserRole")
	}

	// Remover o UserRole
	err = p.permissionRepository.DeleteUserRole(ctx, nil, userRole.GetID())
	if err != nil {
		slog.Error("permission.role.delete_user_role_failed", "user_id", userID, "role_id", roleID, "user_role_id", userRole.GetID(), "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	slog.Info("permission.role.removed", "user_id", userID, "role_id", roleID)
	return nil
}

// RemoveRoleFromUserWithTx remove um role de um usuário usando uma transação existente
func (p *permissionServiceImpl) RemoveRoleFromUserWithTx(ctx context.Context, tx *sql.Tx, userID, roleID int64) error {
	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	if roleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	slog.Debug("permission.role.remove.start.tx", "user_id", userID, "role_id", roleID)

	// Buscar o UserRole existente
	userRole, err := p.permissionRepository.GetUserRoleByUserIDAndRoleID(ctx, tx, userID, roleID)
	if err != nil {
		slog.Error("permission.role.get_user_role_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	if userRole == nil {
		return utils.NotFoundError("UserRole")
	}

	// Remover o UserRole
	err = p.permissionRepository.DeleteUserRole(ctx, tx, userRole.GetID())
	if err != nil {
		slog.Error("permission.role.delete_user_role_failed", "user_id", userID, "role_id", roleID, "user_role_id", userRole.GetID(), "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	slog.Info("permission.role.removed", "user_id", userID, "role_id", roleID)
	return nil
}
