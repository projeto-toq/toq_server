package userservices

import (
	"context"
	"database/sql"

	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetUserByID returns the domain user with the active role eagerly loaded.
// It enforces the invariant that every user must have exactly one valid active role.
// A read-only transaction is used to ensure a consistent view across user and permission reads.
func (us *userService) GetUserByID(ctx context.Context, id int64) (user usermodel.UserInterface, err error) {
	ctx, spanEnd, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return nil, utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	// Iniciar transação somente leitura para leitura consistente
	tx, txErr := us.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		utils.LoggerFromContext(ctx).Error("user.get_by_id.tx_start_error", "error", txErr)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				utils.LoggerFromContext(ctx).Error("user.get_by_id.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	user, err = us.GetUserByIDWithTx(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		utils.LoggerFromContext(ctx).Error("user.get_by_id.tx_commit_error", "error", cmErr)
		return nil, utils.InternalError("Failed to commit transaction")
	}

	return user, nil
}

// GetUserByIDWithTx loads the user and its active role using the provided transaction.
// It enforces the same invariant regarding the active role.
// Português: use esta variante quando já estiver dentro de uma transação; nunca retorne usuário sem active role.
func (us *userService) GetUserByIDWithTx(ctx context.Context, tx *sql.Tx, id int64) (user usermodel.UserInterface, err error) {
	ctx = utils.ContextWithLogger(ctx)

	// Carrega o usuário básico
	user, err = us.repo.GetUserByID(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("User")
		}
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.get_by_id.read_user_error", "error", err, "user_id", id)
		return nil, utils.InternalError("Failed to get user by ID")
	}

	// Carrega a active role via permission service
	activeRole, aerr := us.permissionService.GetActiveUserRoleWithTx(ctx, tx, id)
	if aerr != nil {
		utils.SetSpanError(ctx, aerr)
		utils.LoggerFromContext(ctx).Error("user.get_by_id.read_active_role_error", "error", aerr, "user_id", id)
		return nil, utils.InternalError("Failed to get active user role")
	}

	if activeRole == nil {
		// Invariável do domínio: todo usuário deve ter exatamente um active role válido
		utils.LoggerFromContext(ctx).Error("user.active_role.missing", "user_id", id)
		derr := utils.InternalError("User active role missing")
		utils.SetSpanError(ctx, derr)
		return nil, derr
	}

	user.SetActiveRole(activeRole)
	return user, nil
}
