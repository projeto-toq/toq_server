package userservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// PushOptOut revoga consentimento de notificações push: limpa tokens e seta opt_status=0.
func (us *userService) PushOptOut(ctx context.Context, userID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.push_optout.tx_start_error", "error", err)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				utils.LoggerFromContext(ctx).Error("user.push_optout.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	err = us.pushOptOut(ctx, tx, userID)
	if err != nil {
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.push_optout.tx_commit_error", "error", err)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) pushOptOut(ctx context.Context, tx *sql.Tx, userID int64) (err error) {
	ctx = utils.ContextWithLogger(ctx)
	// Remove all stored device tokens for this user via repository (best-effort)
	if err = us.deviceTokenRepo.RemoveAllByUserID(userID); err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.push_optout.remove_tokens_error", "error", err, "user_id", userID)
		return
	}

	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.push_optout.read_user_error", "error", err, "user_id", userID)
		return
	}
	user.SetOptStatus(false)
	if err = us.repo.UpdateUserByID(ctx, tx, user); err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.push_optout.update_user_error", "error", err, "user_id", userID)
		return
	}

	if err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Usuário rejeitou receber notificações"); err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.push_optout.audit_error", "error", err, "user_id", userID)
		return
	}
	return
}
