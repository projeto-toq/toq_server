package userservices

import (
	"context"
	"database/sql"
	"errors"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// RemoveRoleFromUser remove um role de um usuário
func (us *userService) RemoveRoleFromUser(ctx context.Context, userID, roleID int64) error {
	// Tracing da operação
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	if roleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	logger.Debug("permission.role.remove.start", "user_id", userID, "role_id", roleID)

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.role.remove.tx_start_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.role.remove.tx_rollback_failed", "user_id", userID, "role_id", roleID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	if err = us.RemoveRoleFromUserWithTx(ctx, tx, userID, roleID); err != nil {
		return err
	}

	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.role.remove.tx_commit_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	return nil
}

// RemoveRoleFromUserWithTx remove um role de um usuário usando uma transação existente
//
// Esta função orquestra a remoção de um role com as seguintes etapas:
//  1. Valida parâmetros de entrada (userID, roleID)
//  2. Busca o UserRole existente no banco
//  3. Remove o registro UserRole via repositório
//  4. Invalida cache de permissões do usuário (best-effort)
//
// A invalidação de cache ocorre APÓS a remoção e NÃO bloqueia a operação
// mesmo em caso de falha, garantindo que o fluxo principal seja completado.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (must not be nil)
//   - userID: ID do usuário que terá o role removido (must be > 0)
//   - roleID: ID do role a ser removido (must be > 0)
//
// Returns:
//   - error: Domain error (400/404) ou infrastructure error (500)
//
// Business Rules:
//   - UserID e RoleID devem ser > 0
//   - UserRole deve existir no banco (404 se não encontrado)
//   - Remoção é hard delete (registro é removido fisicamente)
//
// Side Effects:
//   - Remove registro de user_roles table
//   - Invalida cache de permissões do usuário (best-effort)
//   - Registra log Info confirmando remoção
//
// Example:
//
//	if err := us.RemoveRoleFromUserWithTx(ctx, tx, userID, roleID); err != nil {
//	    return err
//	}
func (us *userService) RemoveRoleFromUserWithTx(ctx context.Context, tx *sql.Tx, userID, roleID int64) error {
	// Ensure logger propagation with request_id and trace_id
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Validate input parameters (business rules)
	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	if roleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	logger.Debug("permission.role.remove.start.tx", "user_id", userID, "role_id", roleID)

	// Fetch existing UserRole (domain validation: must exist)
	userRole, err := us.repo.GetUserRoleByUserIDAndRoleID(ctx, tx, userID, roleID)
	if err != nil {
		logger.Error("permission.role.get_user_role_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	if userRole == nil {
		return utils.NotFoundError("UserRole")
	}

	// Perform hard delete from database
	err = us.repo.DeleteUserRole(ctx, tx, userRole.GetID())
	if err != nil {
		// Handle sql.ErrNoRows as success: happens when DELETE affects 0 rows
		// (userRole was loaded in same transaction, so it exists, but may have been already deleted by concurrent operation)
		if !errors.Is(err, sql.ErrNoRows) {
			// Real infrastructure error
			logger.Error("permission.role.delete_user_role_failed", "user_id", userID, "role_id", roleID, "user_role_id", userRole.GetID(), "error", err)
			utils.SetSpanError(ctx, err)
			return utils.InternalError("Failed to delete user role")
		}
		// No rows affected = already deleted = idempotent success, continue
	}

	// Log success (domain event)
	logger.Info("permission.role.removed", "user_id", userID, "role_id", roleID)

	// Invalidate user permissions cache (best-effort, post-commit operation)
	// Failure does not block the operation
	us.permissionService.InvalidateUserCacheSafe(ctx, userID, "remove_role")

	return nil
}
