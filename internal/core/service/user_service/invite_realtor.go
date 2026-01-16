package userservices

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
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
		tokens := realtor.GetDeviceTokens()
		if len(tokens) > 0 {
			// Send to the first available device token
			pushRequest := globalservice.NotificationRequest{
				Type:    globalservice.NotificationTypeFCM,
				Token:   tokens[0].Token,
				Subject: "Nova Proposta de Trabalho",
				Body:    fmt.Sprintf("A imobiliária %s quer trabalhar com você!", agency.GetNickName()),
			}

			err = notificationService.SendNotification(ctx, pushRequest)
			if err != nil {
				utils.SetSpanError(ctx, err)
				logger.Error("user.invite_realtor.send_push_error", "error", err)
				return err
			}
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

		realtorName := phoneNumber
		if realtor != nil {
			realtorName = realtor.GetNickName()
		}
		auditRecord := auditservice.BuildRecordFromContext(
			ctx,
			agency.GetID(),
			auditmodel.AuditTarget{Type: auditmodel.TargetAgencyInvite, ID: invite.GetID()},
			auditmodel.OperationUpdate,
			map[string]any{
				"agency_id":    agency.GetID(),
				"phone":        phoneNumber,
				"realtor_name": realtorName,
				"action":       "invite_updated",
			},
		)
		if err = us.auditService.RecordChange(ctx, tx, auditRecord); err != nil {
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
		// Handle sql.ErrNoRows as success: happens when MySQL UPDATE finds no changes
		// (realtor was loaded from DB, so realtor exists, just no fields changed)
		if errors.Is(err, sql.ErrNoRows) {
			// No rows affected = no changes needed = success
			err = nil
		} else {
			// Real infrastructure error
			utils.SetSpanError(ctx, err)
			logger.Error("user.invite_realtor.update_realtor_error", "error", err, "realtor_id", realtor.GetID())
			return utils.InternalError("Failed to update realtor")
		}
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
		// Handle sql.ErrNoRows as success: happens when MySQL UPDATE finds no changes
		// (invite was loaded in same transaction, so exists, just no fields changed)
		if !errors.Is(err, sql.ErrNoRows) {
			// Real infrastructure error
			utils.SetSpanError(ctx, err)
			logger.Error("user.invite_realtor.update_invite_error", "error", err, "invite_id", invite.GetID())
			return utils.InternalError("Failed to update invite")
		}
		// No changes needed = success, continue
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
		tokens := realtor.GetDeviceTokens()
		if len(tokens) > 0 {
			// Send to the first available device token
			pushRequest := globalservice.NotificationRequest{
				Type:    globalservice.NotificationTypeFCM,
				Token:   tokens[0].Token,
				Subject: "Nova Proposta de Trabalho",
				Body:    fmt.Sprintf("A imobiliária %s quer trabalhar com você!", agency.GetNickName()),
			}

			err = notificationService.SendNotification(ctx, pushRequest)
			if err != nil {
				utils.SetSpanError(ctx, err)
				logger.Error("user.invite_realtor.send_push_error", "error", err)
				return err
			}
		}
	} else {
		err = us.sendSMStoNewRealtor(ctx, phoneNumber, agency)
		if err != nil {
			return
		}
	}

	realtorID := int64(0)
	if realtor != nil {
		realtorID = realtor.GetID()
	}
	auditRecord := auditservice.BuildRecordFromContext(
		ctx,
		agency.GetID(),
		auditmodel.AuditTarget{Type: auditmodel.TargetAgencyInvite, ID: 0},
		auditmodel.OperationCreate,
		map[string]any{
			"agency_id":  agency.GetID(),
			"realtor_id": realtorID,
			"phone":      phoneNumber,
			"push_sent":  push,
			"action":     "invite_created",
		},
	)
	if err = us.auditService.RecordChange(ctx, tx, auditRecord); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.invite_realtor.audit_error", "error", err)
		return
	}

	return
}
