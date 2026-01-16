package userservices

import (
	"context"
	"database/sql"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateOptStatus consolidates opt-in/out behavior with audit and transactions
func (us *userService) UpdateOptStatus(ctx context.Context, opt bool) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil {
		return utils.AuthenticationError("")
	}

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.update_opt_status.tx_start_error", "error", err)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				utils.LoggerFromContext(ctx).Error("user.update_opt_status.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	if err = us.updateOptStatus(ctx, tx, userID, opt); err != nil {
		return
	}

	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.update_opt_status.tx_commit_error", "error", err)
		return utils.InternalError("Failed to commit transaction")
	}
	return
}

func (us *userService) updateOptStatus(ctx context.Context, tx *sql.Tx, userID int64, opt bool) (err error) {
	ctx = utils.ContextWithLogger(ctx)
	// Busca usuário atual para aplicar idempotência
	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.update_opt_status.read_user_error", "error", err, "user_id", userID)
		return
	}

	// Se já está no estado desejado, não faz nada (idempotente)
	if user.IsOptStatus() == opt {
		utils.LoggerFromContext(ctx).Info("user.update_opt_status.idempotent", "user_id", userID, "opt", opt)
		return nil
	}

	// Transição para opt-out: remover tokens antes de persistir
	tokensRemoved := false
	if !opt {
		if err = us.repo.RemoveAllDeviceTokensByUserID(ctx, tx, userID); err != nil {
			utils.LoggerFromContext(ctx).Warn("user.update_opt_status.device_tokens_delete_failed", "user_id", userID, "error", err)
			return
		}
		utils.LoggerFromContext(ctx).Info("user.update_opt_status.device_tokens_deleted", "user_id", userID)
		tokensRemoved = true
	}

	// Atualiza status e persiste
	user.SetOptStatus(opt)
	if err = us.repo.UpdateUserByID(ctx, tx, user); err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.update_opt_status.update_user_error", "error", err, "user_id", userID)
		return
	}

	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		userID,
		auditmodel.AuditTarget{Type: auditmodel.TargetUser, ID: userID},
		auditmodel.OperationUpdate,
		map[string]any{
			"opt_status":            opt,
			"device_tokens_removed": tokensRemoved,
			"trigger":               "update_opt_status",
			"idempotent_before":     user.IsOptStatus(),
		},
	)
	if err = us.auditService.RecordChange(ctx, tx, auditRecord); err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.update_opt_status.audit_error", "error", err, "user_id", userID, "opt", opt)
		return
	}
	return
}
