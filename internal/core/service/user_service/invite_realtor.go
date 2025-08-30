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

func (us *userService) InviteRealtor(ctx context.Context, phoneNumber string) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	err = us.inviteRealtor(ctx, tx, phoneNumber)
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

func (us *userService) inviteRealtor(ctx context.Context, tx *sql.Tx, phoneNumber string) (err error) {

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	//recovery the agency inviting the realtor
	agency, err := us.repo.GetUserByID(ctx, tx, infos.ID)
	if err != nil {
		return
	}

	//verify if the realtor is already on the platform and linked to another agency
	//return error is is already linked
	realtor, isOnPlataform, err := us.isOnPlataform(ctx, tx, phoneNumber)
	if err != nil {
		return
	}

	//verify if the realtor is already invited
	invite, isInvited, err := us.isInvited(ctx, tx, phoneNumber)
	if err != nil {
		return
	}

	switch {
	case isOnPlataform && isInvited:
		err = us.updateInvite(ctx, tx, invite, agency)
		if err != nil {
			return
		}

		// Enviar notificação push para corretor já na plataforma
		notificationService := us.globalService.GetUnifiedNotificationService()
		pushRequest := globalservice.NotificationRequest{
			Type:    globalservice.NotificationTypeFCM,
			Token:   realtor.GetDeviceToken(),
			Subject: "Nova Proposta de Trabalho",
			Body:    fmt.Sprintf("A imobiliária %s quer trabalhar com você!", agency.GetNickName()),
		}

		err = notificationService.SendNotification(ctx, pushRequest)
		if err != nil {
			return err
		}

	case isOnPlataform && !isInvited:
		err = us.createAgencyInviteByPhone(ctx, tx, agency, phoneNumber, realtor, true)
		if err != nil {
			return
		}
	case !isOnPlataform && isInvited:
		err = us.updateInvite(ctx, tx, invite, agency)
		if err != nil {
			return
		}

		err = us.sendSMStoNewRealtor(ctx, phoneNumber, agency)
		if err != nil {
			return err
		}

		err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableRealtorAgency,
			fmt.Sprintf("Atualizado o convite de relacionamento com o Corretor %s", realtor.GetNickName()))
		if err != nil {
			return
		}
	case !isOnPlataform && !isInvited:
		err = us.createAgencyInviteByPhone(ctx, tx, agency, phoneNumber, realtor, false)
		if err != nil {
			return
		}
	}

	//update the user status
	err = us.repo.UpdateUserByID(ctx, tx, realtor)
	if err != nil {
		return
	}

	return
}

func (us *userService) isOnPlataform(ctx context.Context, tx *sql.Tx, phoneNumber string) (realtor usermodel.UserInterface, isOnPlataform bool, err error) {

	realtor, err = us.repo.GetUserByPhoneNumber(ctx, tx, phoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
	}
	if realtor != nil {
		//verify if the realtor is already linked to an agency
		err = us.isAlreadyLinked(ctx, tx, realtor.GetID())
		if err != nil {
			return
		}
	}

	return realtor, realtor != nil, nil
}

func (us *userService) isAlreadyLinked(ctx context.Context, tx *sql.Tx, realtorID int64) (err error) {

	_, err = us.repo.GetAgencyOfRealtor(ctx, tx, realtorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		return nil
	}
	return utils.ErrInternalServer
}

func (us *userService) isInvited(ctx context.Context, tx *sql.Tx, phoneNumber string) (invite usermodel.InviteInterface, isInvited bool, err error) {
	invite, err = us.repo.GetInviteByPhoneNumber(ctx, tx, phoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}

		return nil, false, nil
	}

	return invite, true, nil
}

func (us *userService) updateInvite(ctx context.Context, tx *sql.Tx, invite usermodel.InviteInterface, agency usermodel.UserInterface) (err error) {
	invite.SetAgencyID(agency.GetID())
	err = us.repo.UpdateAgencyInviteByID(ctx, tx, invite)
	if err != nil {
		return
	}

	return
}

func (us *userService) sendSMStoNewRealtor(ctx context.Context, phoneNumber string, agency usermodel.UserInterface) (err error) {
	// Enviar SMS para corretor que não está na plataforma
	notificationService := us.globalService.GetUnifiedNotificationService()
	smsRequest := globalservice.NotificationRequest{
		Type: globalservice.NotificationTypeSMS,
		To:   phoneNumber,
		Body: fmt.Sprintf("A imobiliária %s quer trabalhar com você! Baixe o app TOQ e aceite o convite.", agency.GetNickName()),
	}

	err = notificationService.SendNotification(ctx, smsRequest)
	if err != nil {
		return
	}
	return
}

func (us *userService) createAgencyInviteByPhone(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, phoneNumber string, realtor usermodel.UserInterface, push bool) (err error) {
	err = us.repo.CreateAgencyInvite(ctx, tx, agency, phoneNumber)
	if err != nil {
		return
	}
	// Converter RoleInterface para RoleSlug
	roleSlug := permissionmodel.RoleSlug(realtor.GetActiveRole().GetRole().GetSlug())
	_, _, _, err = us.updateUserStatus(ctx, tx, roleSlug, usermodel.ActionFinishedInviteCreated)
	if err != nil {
		return
	}
	// TODO: Implementar atualização de status via permission service
	// realtor.GetActiveRole().SetStatus(status)
	// realtor.GetActiveRole().SetStatusReason(reason)

	if push {
		// Enviar notificação push para corretor na plataforma
		notificationService := us.globalService.GetUnifiedNotificationService()
		pushRequest := globalservice.NotificationRequest{
			Type:    globalservice.NotificationTypeFCM,
			Token:   realtor.GetDeviceToken(),
			Subject: "Nova Proposta de Trabalho",
			Body:    fmt.Sprintf("A imobiliária %s quer trabalhar com você!", agency.GetNickName()),
		}

		err = notificationService.SendNotification(ctx, pushRequest)
		if err != nil {
			return err
		}
	} else {
		err = us.sendSMStoNewRealtor(ctx, phoneNumber, agency)
		if err != nil {
			return
		}
	}

	err = us.globalService.CreateAudit(ctx, tx, globalmodel.TableRealtorAgency,
		fmt.Sprintf("Criado o convite de relacionamento com o Corretor %s", realtor.GetNickName()))
	if err != nil {
		return
	}

	return
}
