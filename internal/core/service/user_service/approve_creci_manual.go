package userservices

import (
	"context"
	"database/sql"
	"errors"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ApproveCreciManual updates realtor status from pending manual to approved/refused and dispatches FCM notifications to all opted-in devices
func (us *userService) ApproveCreciManual(ctx context.Context, userID int64, target globalmodel.UserRoleStatus) (err error) {
	ctx, spanEnd, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Validate target status: allowed set is Active or one of the refused statuses
	if !globalmodel.IsManualApprovalTarget(target) {
		return utils.ValidationError("status", "Invalid target status")
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("admin.approve_creci.tx_start_error", "error", txErr)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("admin.approve_creci.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	// Check current status must be pending manual
	activeRole, aerr := us.GetActiveUserRoleWithTx(ctx, tx, userID)
	if aerr != nil {
		utils.SetSpanError(ctx, aerr)
		logger.Error("admin.approve_creci.get_active_role_error", "user_id", userID, "error", aerr)
		return utils.InternalError("Failed to get active role")
	}
	if activeRole == nil || activeRole.GetRole() == nil {
		return utils.InternalError("User active role missing")
	}
	if activeRole.GetStatus() != globalmodel.StatusPendingManual {
		return utils.ConflictError("User is not pending manual review")
	}

	// Update status on active role of realtor
	roleSlug := permissionmodel.RoleSlug(activeRole.GetRole().GetSlug())
	statusFrom := activeRole.GetStatus()
	if err = us.repo.UpdateUserRoleStatus(ctx, tx, userID, roleSlug, target); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("Active role to update")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("admin.approve_creci.update_status_error", "user_id", userID, "role", string(roleSlug), "status", target, "error", err)
		return utils.InternalError("Failed to update user role status")
	}

	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		userID,
		auditmodel.AuditTarget{Type: auditmodel.TargetUserRole, ID: activeRole.GetID()},
		auditmodel.OperationStatusChange,
		map[string]any{
			"role_slug":   string(roleSlug),
			"status_from": statusFrom.String(),
			"status_to":   target.String(),
			"reason":      "manual_review",
		},
	)
	if aerr := us.auditService.RecordChange(ctx, tx, auditRecord); aerr != nil {
		utils.SetSpanError(ctx, aerr)
		logger.Error("admin.approve_creci.audit_error", "error", aerr)
		return utils.InternalError("Failed to write audit log")
	}

	if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("admin.approve_creci.tx_commit_error", "error", cmErr)
		return utils.InternalError("Failed to commit transaction")
	}

	if err := us.sendManualApprovalNotification(ctx, userID, target); err != nil {
		return err
	}

	return nil
}

func (us *userService) sendManualApprovalNotification(ctx context.Context, userID int64, target globalmodel.UserRoleStatus) error {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tokens, err := us.globalService.ListDeviceTokensByUserIDIfOptedIn(ctx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("admin.approve_creci.list_tokens_error", "user_id", userID, "error", err)
		return utils.InternalError("Failed to load device tokens")
	}

	if len(tokens) == 0 {
		logger.Info("admin.approve_creci.notification_skipped_no_tokens", "user_id", userID, "target_status", target.String())
		return nil
	}

	subject, body := buildManualApprovalNotificationPayload(target)
	notify := us.globalService.GetUnifiedNotificationService()

	for _, token := range tokens {
		req := globalservice.NotificationRequest{
			Type:    globalservice.NotificationTypeFCM,
			Token:   token,
			Subject: subject,
			Body:    body,
		}

		if err := notify.SendNotification(ctx, req); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("admin.approve_creci.notification_error", "user_id", userID, "token", token, "error", err)
			return utils.InternalError("Failed to dispatch notification")
		}
	}

	logger.Info("admin.approve_creci.notification_sent", "user_id", userID, "target_status", target.String(), "tokens_count", len(tokens))
	return nil
}

func buildManualApprovalNotificationPayload(target globalmodel.UserRoleStatus) (subject, body string) {
	if target == globalmodel.StatusActive {
		return "Aprovação de Cadastro", "Seu cadastro como corretor foi aprovado."
	}

	subject = "Reprovação de Cadastro"
	body = "Seu cadastro foi reprovado."

	switch target {
	case globalmodel.StatusRefusedImage:
		body = "Seu cadastro foi reprovado por problemas nas imagens enviadas."
	case globalmodel.StatusRefusedDocument:
		body = "Seu cadastro foi reprovado por inconsistência nos documentos."
	case globalmodel.StatusRefusedData:
		body = "Seu cadastro foi reprovado por divergência de dados."
	}

	return
}
