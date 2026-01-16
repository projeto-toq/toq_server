package userservices

import (
	"context"
	"database/sql"
	"fmt"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (us *userService) DeleteRealtorOfAgency(ctx context.Context, agencyID int64, realtorID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("user.delete_realtor_of_agency.tracer_error", "error", err)
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_realtor_of_agency.tx_start_error", "error", err)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("user.delete_realtor_of_agency.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	err = us.deleteRealtorOfAgency(ctx, tx, agencyID, realtorID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_realtor_of_agency.delete_relation_error", "error", err, "agency_id", agencyID, "realtor_id", realtorID)
		return err
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_realtor_of_agency.tx_commit_error", "error", err)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) deleteRealtorOfAgency(ctx context.Context, tx *sql.Tx, agencyID int64, realtorID int64) (err error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	realtor, err := us.repo.GetUserByID(ctx, tx, realtorID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_realtor_of_agency.get_realtor_error", "error", err, "realtor_id", realtorID)
		return
	}

	agency, err := us.repo.GetUserByID(ctx, tx, agencyID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_realtor_of_agency.get_agency_error", "error", err, "agency_id", agencyID)
		return
	}

	_, err = us.repo.DeleteAgencyRealtorRelation(ctx, tx, agencyID, realtorID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_realtor_of_agency.delete_relation_error", "error", err, "agency_id", agencyID, "realtor_id", realtorID)
		return
	}

	// Notificar o corretor sobre a remoção da imobiliária
	notificationService := us.globalService.GetUnifiedNotificationService()
	emailRequest := globalservice.NotificationRequest{
		Type:    globalservice.NotificationTypeEmail,
		To:      realtor.GetEmail(),
		Subject: "Remoção de Imobiliária - TOQ",
		Body:    fmt.Sprintf("Você foi removido da imobiliária %s.", agency.GetNickName()),
	}

	err = notificationService.SendNotification(ctx, emailRequest)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.delete_realtor_of_agency.send_notification_error", "error", err, "realtor_id", realtor.GetID(), "agency_id", agency.GetID())
		return
	}

	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		agency.GetID(),
		auditmodel.AuditTarget{Type: auditmodel.TargetRealtorAgency, ID: 0},
		auditmodel.OperationDelete,
		map[string]any{
			"action":       "realtor_removed_by_agency",
			"agency_id":    agency.GetID(),
			"realtor_id":   realtor.GetID(),
			"realtor_name": realtor.GetNickName(),
		},
	)
	if auditErr := us.auditService.RecordChange(ctx, tx, auditRecord); auditErr != nil {
		utils.SetSpanError(ctx, auditErr)
		logger.Error("user.delete_realtor_of_agency.audit_error", "error", auditErr, "realtor_id", realtor.GetID(), "agency_id", agency.GetID())
		return auditErr
	}

	return
}
