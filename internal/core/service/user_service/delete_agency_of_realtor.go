package userservices

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) DeleteAgencyOfRealtor(ctx context.Context, realtorID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("user.delete_agency_of_realtor.tx_start_error", "err", err)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.delete_agency_of_realtor.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	err = us.deleteAgencyOfRealtor(ctx, tx, realtorID)
	if err != nil {
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		slog.Error("user.delete_agency_of_realtor.tx_commit_error", "err", err)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) deleteAgencyOfRealtor(ctx context.Context, tx *sql.Tx, realtorID int64) (err error) {

	agency, err := us.repo.GetAgencyOfRealtor(ctx, tx, realtorID)
	if err != nil {
		return
	}

	realtor, err := us.repo.GetUserByID(ctx, tx, realtorID)
	if err != nil {
		return
	}

	_, err = us.repo.DeleteAgencyRealtorRelation(ctx, tx, agency.GetID(), realtor.GetID())
	if err != nil {
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
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableRealtorAgency,
		fmt.Sprintf("Apagado o relacionamento com a Imobiliária %s", agency.GetNickName()))
	if err != nil {
		return
	}

	return
}
