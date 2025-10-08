package userservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// RequestEmailChange starts the email change flow by generating a validation code
// and persisting the new email as pending. If there is already a pending email
// change (valid or expired), this request regenerates a new code and expiration
// and overwrites the pending entry. The user ID is read from context (SSOT).
func (us *userService) RequestEmailChange(ctx context.Context, newEmail string) (err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	// normalizar email
	newEmail = strings.TrimSpace(strings.ToLower(newEmail))

	// Obter o ID do usuário do contexto
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return utils.InternalError("Failed to resolve user from context")
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		utils.LoggerFromContext(ctx).Error("email_change.request.tx_start_error", "error", txErr)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				utils.LoggerFromContext(ctx).Error("email_change.request.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	user, validation, err := us.requestEmailChange(ctx, tx, userID, newEmail)
	if err != nil {
		return
	}

	// Commit the transaction BEFORE sending notification
	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("email_change.request.tx_commit_error", "error", err)
		return utils.InternalError("Failed to commit transaction")
	}

	// Se não houve pendência criada (mesmo e-mail do atual), retornar sucesso sem notificar
	if validation == nil || validation.GetEmailCode() == "" {
		return nil
	}

	// Enviar notificação (assíncrono pelo serviço unificado)
	notificationService := us.globalService.GetUnifiedNotificationService()
	emailRequest := globalservice.NotificationRequest{
		Type:    globalservice.NotificationTypeEmail,
		To:      validation.GetNewEmail(),
		Subject: "TOQ - Confirmação de Alteração de Email",
		Body:    "Seu código de validação para alteração de email é: " + validation.GetEmailCode(),
	}

	notifyErr := notificationService.SendNotification(ctx, emailRequest)
	if notifyErr != nil {
		utils.SetSpanError(ctx, notifyErr)
		utils.LoggerFromContext(ctx).Error("email_change.request.notification_error", "user_id", user.GetID(), "error", notifyErr)
	}

	return
}

func (us *userService) requestEmailChange(ctx context.Context, tx *sql.Tx, id int64, email string) (user usermodel.UserInterface, validation usermodel.ValidationInterface, err error) {

	user, err = us.repo.GetUserByID(ctx, tx, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("email_change.request.read_user_error", "error", err, "user_id", id)
		return
	}

	// Não tratar mais como no-op quando o novo email é igual ao atual;
	// seguirá como troca comum (sempre (re)gerar código e expiração)

	// Verificar unicidade global (outros usuários não podem ter este email)
	if exist, verr := us.repo.ExistsEmailForAnotherUser(ctx, tx, email, user.GetID()); verr != nil {
		utils.SetSpanError(ctx, verr)
		utils.LoggerFromContext(ctx).Error("email_change.request.exists_email_error", "error", verr, "user_id", user.GetID())
		return nil, nil, verr
	} else if exist {
		return nil, nil, utils.ErrEmailAlreadyInUse
	}

	//set the user validation as pending for email (Option A: sempre sobrescrever com novo código/expiração)
	validation, err = us.repo.GetUserValidations(ctx, tx, user.GetID())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.SetSpanError(ctx, err)
			utils.LoggerFromContext(ctx).Error("email_change.request.read_validations_error", "error", err, "user_id", user.GetID())
			return
		}
		validation = usermodel.NewValidation()
	}
	validation.SetUserID(user.GetID())
	validation.SetEmailCode(us.random6Digits())
	validation.SetEmailCodeExp(time.Now().UTC().Add(usermodel.ValidationCodeExpiration))
	validation.SetNewEmail(email)

	err = us.repo.UpdateUserValidations(ctx, tx, validation)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("email_change.request.update_validations_error", "error", err, "user_id", user.GetID())
		return
	}

	// Note: SendNotification moved to after transaction commit
	// Note: Last activity is now tracked automatically by AuthInterceptor → Redis → Batch worker

	return
}
