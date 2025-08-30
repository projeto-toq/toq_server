package userservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) RejectInvitation(ctx context.Context, realtorID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	err = us.rejectInvitation(ctx, tx, realtorID)
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

func (us *userService) rejectInvitation(ctx context.Context, tx *sql.Tx, userID int64) (err error) {
	//recover the realtor
	realtor, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	//recover the agency invitation
	invite, err := us.repo.GetInviteByPhoneNumber(ctx, tx, realtor.GetPhoneNumber())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrInternalServer
		}
	}

	//recover the agency inviting the realtor
	agency, err := us.repo.GetUserByID(ctx, tx, invite.GetAgencyID())
	if err != nil {
		return
	}

	//delete the invitation
	_, err = us.repo.DeleteInviteByID(ctx, tx, invite.GetID())
	if err != nil {
		return
	}

	// Converter RoleInterface para RoleSlug
	roleSlug := permissionmodel.RoleSlug(realtor.GetActiveRole().GetRole().GetSlug())
	_, _, _, err = us.updateUserStatus(ctx, tx, roleSlug, usermodel.ActionFinishedInviteRejected)
	if err != nil {
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
		return
	}

	err = us.repo.UpdateUserByID(ctx, tx, realtor)
	if err != nil {
		return
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableAgencyInvites, "Convite para relacionamento com imobiliária rejeitado")
	if err != nil {
		return
	}

	return
}
