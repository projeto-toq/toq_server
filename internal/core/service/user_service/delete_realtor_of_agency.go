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

func (us *userService) DeleteRealtorOfAgency(ctx context.Context, agencyID int64, realtorID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		slog.Error("user.delete_realtor_of_agency.tx_start_error", "err", err)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.delete_realtor_of_agency.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	err = us.deleteRealtorOfAgency(ctx, tx, agencyID, realtorID)
	if err != nil {
		return err
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		slog.Error("user.delete_realtor_of_agency.tx_commit_error", "err", err)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) deleteRealtorOfAgency(ctx context.Context, tx *sql.Tx, agencyID int64, realtorID int64) (err error) {

	realtor, err := us.repo.GetUserByID(ctx, tx, realtorID)
	if err != nil {
		return
	}

	agency, err := us.repo.GetUserByID(ctx, tx, agencyID)
	if err != nil {
		return
	}

	_, err = us.repo.DeleteAgencyRealtorRelation(ctx, tx, agencyID, realtorID)
	if err != nil {
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
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableRealtorAgency,
		fmt.Sprintf("Apagado o relacionamento com o Corretor %s", realtor.GetNickName()))
	if err != nil {
		return
	}

	return
}
