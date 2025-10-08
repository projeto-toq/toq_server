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
	validators "github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

func (us *userService) InviteRealtor(ctx context.Context, phoneNumber string) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if phoneNumber != "" {
		normalizedPhone, normErr := validators.NormalizeToE164(phoneNumber)
		if normErr != nil {
			return normErr
		}
		phoneNumber = normalizedPhone
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("user.invite_realtor.tx_start_error", "error", txErr)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("user.invite_realtor.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	err = us.inviteRealtor(ctx, tx, phoneNumber)
	if err != nil {
		return
	}

	if cmErr := us.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("user.invite_realtor.tx_commit_error", "error", cmErr)
		return utils.InternalError("Failed to commit transaction")
	}

	return
}

func (us *userService) inviteRealtor(ctx context.Context, tx *sql.Tx, phoneNumber string) (err error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	//recovery the agency inviting the realtor
	agency, err := us.repo.GetUserByID(ctx, tx, infos.ID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.invite_realtor.read_agency_error", "error", err, "agency_id", infos.ID)
		return
	}

	//verify if the realtor is already on the platform and linked to another agency
	//return error is is already linked
	realtor, isOnPlataform, err := us.isOnPlataform(ctx, tx, phoneNumber)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.invite_realtor.check_platform_error", "error", err, "phone", phoneNumber)
		return
	}

	//verify if the realtor is already invited
	invite, isInvited, err := us.isInvited(ctx, tx, phoneNumber)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.invite_realtor.check_invited_error", "error", err, "phone", phoneNumber)
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
			utils.SetSpanError(ctx, err)
			logger.Error("user.invite_realtor.send_push_error", "error", err)
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
		utils.SetSpanError(ctx, err)
		logger.Error("user.invite_realtor.update_realtor_error", "error", err, "realtor_id", realtor.GetID())
		return
	}

	return
}

func (us *userService) isOnPlataform(ctx context.Context, tx *sql.Tx, phoneNumber string) (realtor usermodel.UserInterface, isOnPlataform bool, err error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	realtor, err = us.repo.GetUserByPhoneNumber(ctx, tx, phoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		utils.SetSpanError(ctx, err)
		logger.Error("user.invite_realtor.read_realtor_by_phone_error", "error", err, "phone", phoneNumber)
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
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	_, err = us.repo.GetAgencyOfRealtor(ctx, tx, realtorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		utils.SetSpanError(ctx, err)
		logger.Error("user.invite_realtor.read_agency_of_realtor_error", "error", err, "realtor_id", realtorID)
		return nil
	}
	return utils.ConflictError("Realtor already linked to an agency")
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
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	invite.SetAgencyID(agency.GetID())
	err = us.repo.UpdateAgencyInviteByID(ctx, tx, invite)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.invite_realtor.update_invite_error", "error", err, "invite_id", invite.GetID())
		return
	}

	return
}

func (us *userService) sendSMStoNewRealtor(ctx context.Context, phoneNumber string, agency usermodel.UserInterface) (err error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	// Enviar SMS para corretor que não está na plataforma
	notificationService := us.globalService.GetUnifiedNotificationService()
	smsRequest := globalservice.NotificationRequest{
		Type: globalservice.NotificationTypeSMS,
		To:   phoneNumber,
		Body: fmt.Sprintf("A imobiliária %s quer trabalhar com você! Baixe o app TOQ e aceite o convite.", agency.GetNickName()),
	}

	err = notificationService.SendNotification(ctx, smsRequest)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.invite_realtor.send_sms_error", "error", err)
		return
	}
	return
}

func (us *userService) createAgencyInviteByPhone(ctx context.Context, tx *sql.Tx, agency usermodel.UserInterface, phoneNumber string, realtor usermodel.UserInterface, push bool) (err error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	err = us.repo.CreateAgencyInvite(ctx, tx, agency, phoneNumber)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.invite_realtor.create_invite_error", "error", err)
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
			utils.SetSpanError(ctx, err)
			logger.Error("user.invite_realtor.send_push_error", "error", err)
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
		utils.SetSpanError(ctx, err)
		logger.Error("user.invite_realtor.audit_error", "error", err)
		return
	}

	return
}
