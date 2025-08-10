package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// PushOptOut revoga consentimento de notificações push: limpa tokens e seta opt_status=0.
func (us *userService) PushOptOut(ctx context.Context, userID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	err = us.pushOptOut(ctx, tx, userID)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	return
}

func (us *userService) pushOptOut(ctx context.Context, tx *sql.Tx, userID int64) (err error) {
	// Remove all stored device tokens for this user via repository (best-effort)
	if err = us.repo.RemoveAllDeviceTokens(ctx, tx, userID); err != nil {
		return
	}

	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}
	user.SetOptStatus(false)
	if err = us.repo.UpdateUserByID(ctx, tx, user); err != nil {
		return
	}

	if err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Usuário rejeitou receber notificações"); err != nil {
		return
	}
	return
}
