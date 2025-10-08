package userservices

import (
	"context"
	"database/sql"
	"fmt"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) DeleteAgencyOfRealtor(ctx context.Context, realtorID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		ctx = utils.ContextWithLogger(ctx)
		utils.LoggerFromContext(ctx).Error("user.delete_agency_of_realtor.tracer_error", "error", err)
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.delete_agency_of_realtor.tx_start_error", "error", err)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				utils.LoggerFromContext(ctx).Error("user.delete_agency_of_realtor.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	err = us.deleteAgencyOfRealtor(ctx, tx, realtorID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.delete_agency_of_realtor.delete_relation_error", "error", err, "realtor_id", realtorID)
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.delete_agency_of_realtor.tx_commit_error", "error", err)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) deleteAgencyOfRealtor(ctx context.Context, tx *sql.Tx, realtorID int64) (err error) {
	ctx = utils.ContextWithLogger(ctx)

	agency, err := us.repo.GetAgencyOfRealtor(ctx, tx, realtorID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.delete_agency_of_realtor.get_agency_error", "error", err, "realtor_id", realtorID)
		return
	}

	realtor, err := us.repo.GetUserByID(ctx, tx, realtorID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.delete_agency_of_realtor.get_realtor_error", "error", err, "realtor_id", realtorID)
		return
	}

	_, err = us.repo.DeleteAgencyRealtorRelation(ctx, tx, agency.GetID(), realtor.GetID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.delete_agency_of_realtor.delete_relation_error", "error", err, "agency_id", agency.GetID(), "realtor_id", realtor.GetID())
		return
	}

	// Notificar o corretor sobre a saída da imobiliária
	notificationService := us.globalService.GetUnifiedNotificationService()
	emailRequest := globalservice.NotificationRequest{
		Type:    globalservice.NotificationTypeEmail,
		To:      realtor.GetEmail(),
		Subject: "Saída de Imobiliária - TOQ",
		Body:    fmt.Sprintf("Você saiu da imobiliária %s.", agency.GetNickName()),
	}

	err = notificationService.SendNotification(ctx, emailRequest)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.delete_agency_of_realtor.send_notification_error", "error", err, "realtor_id", realtor.GetID(), "agency_id", agency.GetID())
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableRealtorAgency,
		fmt.Sprintf("Apagado o relacionamento com a Imobiliária %s", agency.GetNickName()))
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("user.delete_agency_of_realtor.audit_error", "error", err, "realtor_id", realtor.GetID(), "agency_id", agency.GetID())
		return
	}

	return
}
