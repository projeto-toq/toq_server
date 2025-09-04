package userservices

import (
	"context"
	"database/sql"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// PushOptIn sets opt_status=1 for the user indicating consent to receive push notifications.
// Device tokens are managed during SignIn; this only persists the consent flag.
func (us *userService) PushOptIn(ctx context.Context, userID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("user.push_optin.tx_start_error", "err", err)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.push_optin.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	err = us.pushOptIn(ctx, tx, userID)
	if err != nil {
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		slog.Error("user.push_optin.tx_commit_error", "err", err)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) pushOptIn(ctx context.Context, tx *sql.Tx, userID int64) (err error) {
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	user.SetOptStatus(true)

	if err = us.repo.UpdateUserByID(ctx, tx, user); err != nil {
		return
	}

	if err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Usuário aceitou receber notificações"); err != nil {
		return
	}

	return
}
