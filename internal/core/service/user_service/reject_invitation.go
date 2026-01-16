package userservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (us *userService) RejectInvitation(ctx context.Context, realtorID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.reject_invitation.tx_start_error", "error", err)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("user.reject_invitation.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	err = us.rejectInvitation(ctx, tx, realtorID)
	if err != nil {
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.reject_invitation.tx_commit_error", "error", err)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) rejectInvitation(ctx context.Context, tx *sql.Tx, userID int64) (err error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	//recover the realtor
	realtor, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.reject_invitation.read_realtor_error", "error", err, "realtor_id", userID)
		return
	}

	//recover the agency invitation
	invite, err := us.repo.GetInviteByPhoneNumber(ctx, tx, realtor.GetPhoneNumber())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("Invitation")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("user.reject_invitation.read_invite_error", "error", err, "phone", realtor.GetPhoneNumber())
		return utils.InternalError("Failed to get invitation")
	}

	//recover the agency inviting the realtor
	agency, err := us.repo.GetUserByID(ctx, tx, invite.GetAgencyID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.reject_invitation.read_agency_error", "error", err, "agency_id", invite.GetAgencyID())
		return
	}

	//delete the invitation
	_, err = us.repo.DeleteInviteByID(ctx, tx, invite.GetID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.reject_invitation.delete_invite_error", "error", err, "invite_id", invite.GetID())
		return
	}

	// Converter RoleInterface para RoleSlug
	roleSlug := permissionmodel.RoleSlug(realtor.GetActiveRole().GetRole().GetSlug())
	_, _, _, err = us.updateUserStatus(ctx, tx, roleSlug, usermodel.ActionFinishedInviteRejected)
	if err != nil {
		// keep domain/infra mapping inside updateUserStatus
		return
	}
	// TODO: Implementar atualização de status via permission service
	// realtor.GetActiveRole().SetStatus(status)
	// realtor.GetActiveRole().SetStatusReason(reason)

	// Notificar a imobiliária sobre a rejeição do convite
	notificationService := us.globalService.GetUnifiedNotificationService()
	emailRequest := globalservice.NotificationRequest{
		Type:    globalservice.NotificationTypeEmail,
		To:      agency.GetEmail(),
		Subject: "Convite Rejeitado - TOQ",
		Body:    fmt.Sprintf("O corretor %s rejeitou seu convite para trabalhar com sua imobiliária.", realtor.GetFullName()),
	}

	err = notificationService.SendNotification(ctx, emailRequest)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.reject_invitation.send_notification_error", "error", err)
		return err
	}

	err = us.repo.UpdateUserByID(ctx, tx, realtor)
	if err != nil {
		// Handle sql.ErrNoRows as success: happens when MySQL UPDATE finds no changes
		// (realtor was loaded in same transaction, so realtor exists, just no fields changed)
		if !errors.Is(err, sql.ErrNoRows) {
			// Real infrastructure error
			utils.SetSpanError(ctx, err)
			logger.Error("user.reject_invitation.update_realtor_error", "error", err, "realtor_id", realtor.GetID())
			return utils.InternalError("Failed to update realtor")
		}
		// No changes needed = success, continue
	}

	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		realtor.GetID(),
		auditmodel.AuditTarget{Type: auditmodel.TargetAgencyInvite, ID: invite.GetID()},
		auditmodel.OperationStatusChange,
		map[string]any{
			"action":      "invite_rejected",
			"agency_id":   agency.GetID(),
			"realtor_id":  realtor.GetID(),
			"phone":       realtor.GetPhoneNumber(),
			"status_from": "pending",
			"status_to":   "rejected",
		},
	)
	if auditErr := us.auditService.RecordChange(ctx, tx, auditRecord); auditErr != nil {
		utils.SetSpanError(ctx, auditErr)
		logger.Error("user.reject_invitation.audit_error", "error", auditErr)
		return auditErr
	}

	return
}
