package userservices

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// RequestEmailChange starts the email change flow by generating a validation code
// and persisting the new email as pending. If the new email equals the current one,
// the operation is a no-op (no pending created, no notification). The user ID is
// read from context (SSOT).
func (us *userService) RequestEmailChange(ctx context.Context, newEmail string) (err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// normalizar email
	newEmail = strings.TrimSpace(strings.ToLower(newEmail))

	// Obter o ID do usuário do contexto
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return utils.InternalError("Failed to resolve user from context")
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		slog.Error("email_change.request.tx_start_error", "err", txErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementEmailChangeRequest("start_tx_error")
		}
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("email_change.request.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	user, validation, err := us.requestEmailChange(ctx, tx, userID, newEmail)
	if err != nil {
		if mp := us.globalService.GetMetrics(); mp != nil {
			// map some known domain errors for better observability
			switch err {
			case utils.ErrEmailAlreadyInUse:
				mp.IncrementEmailChangeRequest("already_in_use")
			default:
				mp.IncrementEmailChangeRequest("domain_error")
			}
		}
		return
	}

	// Commit the transaction BEFORE sending notification
	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		slog.Error("email_change.request.tx_commit_error", "err", err)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementEmailChangeRequest("commit_error")
		}
		return utils.InternalError("Failed to commit transaction")
	}

	// Se não houve pendência criada (mesmo e-mail do atual), retornar sucesso sem notificar
	if validation == nil || validation.GetEmailCode() == "" {
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementEmailChangeRequest("success_noop")
		}
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
		slog.Error("email_change.request.notification_error", "userID", user.GetID(), "err", notifyErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementEmailChangeRequest("notify_error")
		}
	} else if mp := us.globalService.GetMetrics(); mp != nil {
		mp.IncrementEmailChangeRequest("success")
	}

	return
}

func (us *userService) requestEmailChange(ctx context.Context, tx *sql.Tx, id int64, email string) (user usermodel.UserInterface, validation usermodel.ValidationInterface, err error) {

	user, err = us.repo.GetUserByID(ctx, tx, id)
	if err != nil {
		return
	}

	// No-op: se o novo email for igual ao atual, não criar pendência nem enviar código
	if strings.EqualFold(user.GetEmail(), email) {
		// Comentário: manter comportamento idempotente conforme regra de negócio
		return user, nil, nil
	}

	// Verificar unicidade global (outros usuários não podem ter este email)
	if exist, verr := us.repo.ExistsEmailForAnotherUser(ctx, tx, email, user.GetID()); verr != nil {
		return nil, nil, verr
	} else if exist {
		return nil, nil, utils.ErrEmailAlreadyInUse
	}

	//set the user validation as pending for email (Option A: sempre sobrescrever com novo código/expiração)
	validation, err = us.repo.GetUserValidations(ctx, tx, user.GetID())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
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
		return
	}

	// Note: SendNotification moved to after transaction commit
	// Note: Last activity is now tracked automatically by AuthInterceptor → Redis → Batch worker

	return
}
