package userservices

import (
	"context"
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// SwitchActiveRole desativa todos os roles do usuário e ativa apenas o especificado
func (us *userService) SwitchActiveRole(ctx context.Context, userID, newRoleID int64) error {
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	if newRoleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	logger.Debug("permission.role.switch.request", "user_id", userID, "new_role_id", newRoleID)

	// 1. Desativar todos os roles do usuário
	if err := us.DeactivateAllUserRoles(ctx, userID); err != nil {
		logger.Error("permission.role.switch.deactivate_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return err
	}

	// 2. Ativar o novo role
	if err := us.ActivateUserRole(ctx, userID, newRoleID); err != nil {
		logger.Error("permission.role.switch.activate_failed", "user_id", userID, "new_role_id", newRoleID, "error", err)
		utils.SetSpanError(ctx, err)
		return err
	}

	logger.Info("permission.role.switched", "user_id", userID, "new_role_id", newRoleID)
	return nil
}

// SwitchActiveRoleWithTx desativa todos os roles do usuário e ativa apenas o especificado (com transação)
//
// Esta função orquestra a troca de role ativo com as seguintes etapas:
//  1. Valida parâmetros de entrada (userID, newRoleID)
//  2. Desativa todos os roles do usuário (SET is_active = 0)
//  3. Ativa apenas o novo role especificado (SET is_active = 1)
//  4. Invalida cache de permissões do usuário (best-effort)
//
// A invalidação de cache ocorre APÓS ambas as operações e NÃO bloqueia
// o fluxo principal mesmo em caso de falha.
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - tx: Database transaction (must not be nil)
//   - userID: ID do usuário cujo role será trocado (must be > 0)
//   - newRoleID: ID do novo role a ser ativado (must be > 0)
//
// Returns:
//   - error: Domain error (400) ou infrastructure error (500)
//
// Business Rules:
//   - UserID e newRoleID devem ser > 0
//   - Operação é atômica: ambos updates (deactivate all + activate one) dentro da mesma tx
//   - Se o newRoleID não existir, a operação falha (repository retorna erro)
//
// Side Effects:
//   - Atualiza múltiplos registros em user_roles table (is_active = 0)
//   - Ativa um registro específico (is_active = 1)
//   - Invalida cache de permissões do usuário (best-effort)
//   - Registra log Info confirmando a troca
//
// Example:
//
//	if err := us.SwitchActiveRoleWithTx(ctx, tx, userID, newRoleID); err != nil {
//	    return err
//	}
func (us *userService) SwitchActiveRoleWithTx(ctx context.Context, tx *sql.Tx, userID, newRoleID int64) error {
	// Initialize tracing for observability
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	// Ensure logger propagation with request_id and trace_id
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Validate input parameters (business rules)
	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	if newRoleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	logger.Debug("permission.role.switch.tx.request", "user_id", userID, "new_role_id", newRoleID)

	// Step 1: Deactivate all user roles (atomically within transaction)
	if err := us.repo.DeactivateAllUserRoles(ctx, tx, userID); err != nil {
		logger.Error("permission.role.switch.tx.deactivate_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Step 2: Activate only the specified new role (atomically within same transaction)
	if err := us.repo.ActivateUserRole(ctx, tx, userID, newRoleID); err != nil {
		logger.Error("permission.role.switch.tx.activate_failed", "user_id", userID, "new_role_id", newRoleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Log success (domain event)
	logger.Info("permission.role.switched.tx", "user_id", userID, "new_role_id", newRoleID)

	// Invalidate user permissions cache (best-effort, post-operation)
	// Failure does not block the operation as cache will be eventually consistent
	us.permissionService.InvalidateUserCacheSafe(ctx, userID, "switch_active_role")

	return nil
}

// GetActiveUserRole methods moved to get_active_user_role.go

// DeactivateAllUserRoles desativa todos os roles de um usuário
//
// Esta função cria uma transação isolada e:
//  1. Desativa todos os roles do usuário (SET is_active = 0 para todas as linhas)
//  2. Invalida cache de permissões do usuário (best-effort)
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - userID: ID do usuário cujos roles serão desativados (must be > 0)
//
// Returns:
//   - error: Domain error (400) ou infrastructure error (500)
//
// Business Rules:
//   - UserID deve ser > 0
//   - Operação usa transação isolada (não reusa tx de caller)
//
// Side Effects:
//   - Atualiza múltiplos registros em user_roles table
//   - Invalida cache de permissões (best-effort)
//   - Registra log Info confirmando desativação
//
// Example:
//
//	if err := us.DeactivateAllUserRoles(ctx, userID); err != nil {
//	    return err
//	}
func (us *userService) DeactivateAllUserRoles(ctx context.Context, userID int64) error {
	// Initialize tracing for observability
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	// Ensure logger propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Validate input parameter
	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	logger.Debug("permission.user_roles.deactivate.request", "user_id", userID)

	// Start isolated transaction for this operation
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user_roles.deactivate.tx_start_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	defer func() {
		// Rollback on error (explicit commit at the end if success)
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user_roles.deactivate.tx_rollback_failed", "user_id", userID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	// Deactivate all user roles via repository
	if err = us.repo.DeactivateAllUserRoles(ctx, tx, userID); err != nil {
		logger.Error("permission.user_roles.deactivate.db_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Commit the transaction
	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.user_roles.deactivate.tx_commit_failed", "user_id", userID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Log success
	logger.Info("permission.user_roles.deactivated", "user_id", userID)

	// Invalidate cache (best-effort, post-commit)
	us.permissionService.InvalidateUserCacheSafe(ctx, userID, "deactivate_all_roles")

	return nil
}

// ActivateUserRole ativa um role específico do usuário
//
// Esta função cria uma transação isolada e:
//  1. Ativa um role específico (SET is_active = 1)
//  2. Invalida cache de permissões do usuário (best-effort)
//
// Parameters:
//   - ctx: Context for tracing and logging
//   - userID: ID do usuário cujo role será ativado (must be > 0)
//   - roleID: ID do role a ser ativado (must be > 0)
//
// Returns:
//   - error: Domain error (400) ou infrastructure error (500)
//
// Business Rules:
//   - UserID e roleID devem ser > 0
//   - Operação usa transação isolada
//   - Se o roleID não existir para o usuário, repository retorna erro
//
// Side Effects:
//   - Atualiza um registro em user_roles table
//   - Invalida cache de permissões (best-effort)
//   - NÃO desativa outros roles (caller deve fazer isso explicitamente)
//
// Example:
//
//	// Usado em fluxos de troca de role (após desativar todos)
//	if err := us.ActivateUserRole(ctx, userID, roleID); err != nil {
//	    return err
//	}
func (us *userService) ActivateUserRole(ctx context.Context, userID, roleID int64) error {
	// Initialize tracing for observability
	ctx, end, _ := utils.GenerateTracer(ctx)
	defer end()

	// Ensure logger propagation
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Validate input parameters
	if userID <= 0 {
		return utils.BadRequest("invalid user id")
	}

	if roleID <= 0 {
		return utils.BadRequest("invalid role id")
	}

	// Start isolated transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		logger.Error("permission.user_role.activate.tx_start_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				logger.Error("permission.user_role.activate.tx_rollback_failed", "user_id", userID, "role_id", roleID, "error", rbErr)
				utils.SetSpanError(ctx, rbErr)
			}
		}
	}()

	// Activate specific role via repository
	if err = us.repo.ActivateUserRole(ctx, tx, userID, roleID); err != nil {
		logger.Error("permission.user_role.activate.db_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Commit the transaction
	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		logger.Error("permission.user_role.activate.tx_commit_failed", "user_id", userID, "role_id", roleID, "error", err)
		utils.SetSpanError(ctx, err)
		return utils.InternalError("")
	}

	// Invalidate cache (best-effort, post-commit)
	us.permissionService.InvalidateUserCacheSafe(ctx, userID, "activate_user_role")

	return nil
}
