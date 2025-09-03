package userservices

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// ResendEmailChangeCode regenerates the email change code and extends its expiration.
// It requires a pending email change; after commit, sends the new code to the new email address.
// The user ID is extracted from context (SSOT).
func (us *userService) ResendEmailChangeCode(ctx context.Context) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Obter o ID do usuário a partir do contexto (fonte única de verdade)
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return utils.ErrInternalServer
	}

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementEmailChangeResend("start_tx_error")
		}
		return
	}

	var userEmail, newCode string
	userEmail, newCode, err = us.resendEmailChangeCode(ctx, tx, userID)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		if mp := us.globalService.GetMetrics(); mp != nil {
			if errors.Is(err, utils.ErrEmailChangeNotPending) {
				mp.IncrementEmailChangeResend("not_pending")
			} else {
				mp.IncrementEmailChangeResend("domain_error")
			}
		}
		return
	}

	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementEmailChangeResend("commit_error")
		}
		return
	}

	// Após o commit, enviar a notificação por e-mail com o novo código
	notificationService := us.globalService.GetUnifiedNotificationService()
	emailRequest := globalservice.NotificationRequest{
		Type:    globalservice.NotificationTypeEmail,
		To:      userEmail,
		Subject: "TOQ - Novo código de alteração de email",
		Body:    "Seu novo código de validação para alteração de email é: " + newCode,
	}
	if notifyErr := notificationService.SendNotification(ctx, emailRequest); notifyErr != nil {
		slog.Error("Failed to send email notification (resend email change)", "userID", userID, "error", notifyErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementEmailChangeResend("notify_error")
		}
	} else {
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementEmailChangeResend("success")
		}
	}

	return
}

// resendEmailChangeCode performs the regeneration of the email code and extends the expiration.
// Returns the destination email (new email) for notification purposes.
func (us *userService) resendEmailChangeCode(ctx context.Context, tx *sql.Tx, userID int64) (destEmail string, code string, err error) {
	validation, err := us.repo.GetUserValidations(ctx, tx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", utils.ErrEmailChangeNotPending
		}
		return "", "", err
	}

	// Deve existir um novo email pendente
	destEmail = validation.GetNewEmail()
	if destEmail == "" {
		return "", "", utils.ErrEmailChangeNotPending
	}

	// Gerar novo código e estender a expiração
	code = us.random6Digits()
	validation.SetEmailCode(code)
	validation.SetEmailCodeExp(time.Now().UTC().Add(usermodel.ValidationCodeExpiration))

	if err = us.repo.UpdateUserValidations(ctx, tx, validation); err != nil {
		return "", "", err
	}

	return destEmail, code, nil
}
