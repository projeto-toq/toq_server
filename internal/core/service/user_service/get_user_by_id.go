package userservices

import (
	"context"
	"database/sql"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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

	// Iniciar transação somente leitura para leitura consistente
	tx, txErr := us.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		slog.Error("user.get_by_id.tx_start_error", "err", txErr)
		return nil, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.get_by_id.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	user, err = us.GetUserByIDWithTx(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		slog.Error("user.get_by_id.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("Failed to commit transaction")
	}

	return user, nil
}

// GetUserByIDWithTx loads the user and its active role using the provided transaction.
// It enforces the same invariant regarding the active role.
// Português: use esta variante quando já estiver dentro de uma transação; nunca retorne usuário sem active role.
func (us *userService) GetUserByIDWithTx(ctx context.Context, tx *sql.Tx, id int64) (user usermodel.UserInterface, err error) {
	ctx, spanEnd, _ := utils.GenerateTracer(ctx)
	defer spanEnd()

	// Carrega o usuário básico
	user, err = us.repo.GetUserByID(ctx, tx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.NotFoundError("User")
		}
		utils.SetSpanError(ctx, err)
		return nil, utils.InternalError("Failed to get user by ID")
	}

	// Carrega a active role via permission service
	activeRole, aerr := us.permissionService.GetActiveUserRoleWithTx(ctx, tx, id)
	if aerr != nil {
		utils.SetSpanError(ctx, aerr)
		return nil, utils.InternalError("Failed to get active user role")
	}

	if activeRole == nil {
		// Invariável do domínio: todo usuário deve ter exatamente um active role válido
		slog.Error("user.active_role.missing", "user_id", id)
		derr := utils.InternalError("User active role missing")
		utils.SetSpanError(ctx, derr)
		return nil, derr
	}

	user.SetActiveRole(activeRole)
	return user, nil
}
