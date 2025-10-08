package userservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// ApproveCreciManual updates realtor status from pending manual to approved/refused and sends notification
func (us *userService) ApproveCreciManual(ctx context.Context, userID int64, target permissionmodel.UserRoleStatus) (err error) {
	ctx, spanEnd, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Validate target status: allowed set is Active or one of the refused statuses
	if !permissionmodel.IsManualApprovalTarget(target) {
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
	activeRole, aerr := us.permissionService.GetActiveUserRoleWithTx(ctx, tx, userID)
	if aerr != nil {
		utils.SetSpanError(ctx, aerr)
		logger.Error("admin.approve_creci.get_active_role_error", "user_id", userID, "error", aerr)
		return utils.InternalError("Failed to get active role")
	}
	if activeRole == nil || activeRole.GetRole() == nil {
		return utils.InternalError("User active role missing")
	}
	if activeRole.GetStatus() != permissionmodel.StatusPendingManual {
		return utils.ConflictError("User is not pending manual review")
	}

	// Update status on active role of realtor
	roleSlug := permissionmodel.RoleSlug(activeRole.GetRole().GetSlug())
	if err = us.repo.UpdateUserRoleStatus(ctx, tx, userID, roleSlug, target); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("Active role to update")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("admin.approve_creci.update_status_error", "user_id", userID, "role", string(roleSlug), "status", target, "error", err)
		return utils.InternalError("Failed to update user role status")
	}

	if aerr := us.globalService.CreateAudit(ctx, tx, "user_roles", fmt.Sprintf("Admin manual review: set status to %s", target.String())); aerr != nil {
		utils.SetSpanError(ctx, aerr)
		logger.Error("admin.approve_creci.audit_error", "error", aerr)
		return utils.InternalError("Failed to write audit log")
	}

	if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("admin.approve_creci.tx_commit_error", "error", cmErr)
		return utils.InternalError("Failed to commit transaction")
	}

	// After commit: send push notification if device tokens exist
	// Build subject/body by status
	var subject, body string
	if target == permissionmodel.StatusActive {
		subject = "Aprovação de Cadastro"
		body = "Seu cadastro como corretor foi aprovado."
	} else {
		subject = "Reprovação de Cadastro"
		switch target {
		case permissionmodel.StatusRefusedImage:
			body = "Seu cadastro foi reprovado por problemas nas imagens enviadas."
		case permissionmodel.StatusRefusedDocument:
			body = "Seu cadastro foi reprovado por inconsistência nos documentos."
		case permissionmodel.StatusRefusedData:
			body = "Seu cadastro foi reprovado por divergência de dados."
		default:
			body = "Seu cadastro foi reprovado."
		}
	}

	// Load user to obtain last known device token stored during sign-in
	// Using a read-only tx internally
	user, gerr := us.GetUserByID(ctx, userID)
	if gerr == nil && user != nil {
		token := user.GetDeviceToken()
		if token != "" {
			notify := us.globalService.GetUnifiedNotificationService()
			_ = notify.SendNotification(ctx, globalservice.NotificationRequest{
				Type:    globalservice.NotificationTypeFCM,
				Token:   token,
				Subject: subject,
				Body:    body,
			})
		}
	}

	return nil
}
