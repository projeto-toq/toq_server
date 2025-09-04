package userservices

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// ResendEmailChangeCode regenerates the email change code and extends its expiration.
// It requires a pending email change; after commit, sends the new code to the new email address.
// The user ID is extracted from context (SSOT).
func (us *userService) ResendEmailChangeCode(ctx context.Context) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Obter o ID do usuário a partir do contexto (fonte única de verdade)
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return utils.AuthenticationError("")
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		slog.Error("email_change.resend.tx_start_error", "err", txErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementEmailChangeResend("start_tx_error")
		}
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("email_change.resend.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	var userEmail, code string
	userEmail, code, err = us.resendEmailChangeCode(ctx, tx, userID)
	if err != nil {
		if mp := us.globalService.GetMetrics(); mp != nil {
			switch err {
			case utils.ErrEmailChangeNotPending:
				mp.IncrementEmailChangeResend("not_pending")
			case utils.ErrEmailChangeCodeExpired:
				mp.IncrementEmailChangeResend("expired")
			case utils.ErrEmailAlreadyInUse:
				mp.IncrementEmailChangeResend("already_in_use")
			default:
				mp.IncrementEmailChangeResend("domain_error")
			}
		}
		return
	}

	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		slog.Error("email_change.resend.tx_commit_error", "err", err)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementEmailChangeResend("commit_error")
		}
		return utils.InternalError("Failed to commit transaction")
	}

	// Após o commit, enviar a notificação por e-mail com o novo código
	notificationService := us.globalService.GetUnifiedNotificationService()
	emailRequest := globalservice.NotificationRequest{
		Type:    globalservice.NotificationTypeEmail,
		To:      userEmail,
		Subject: "TOQ - Código de alteração de email",
		Body:    "Seu código de validação para alteração de email é: " + code,
	}
	if notifyErr := notificationService.SendNotification(ctx, emailRequest); notifyErr != nil {
		slog.Error("email_change.resend.notification_error", "userID", userID, "err", notifyErr)
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

	// Deve haver um código válido ainda dentro do prazo
	code = validation.GetEmailCode()
	if code == "" {
		return "", "", utils.ErrEmailChangeNotPending
	}
	if validation.GetEmailCodeExp().Before(time.Now().UTC()) {
		return "", "", utils.ErrEmailChangeCodeExpired
	}

	// Verificar unicidade global (outros usuários não podem ter este email)
	if exist, verr := us.repo.ExistsEmailForAnotherUser(ctx, tx, destEmail, userID); verr != nil {
		return "", "", verr
	} else if exist {
		return "", "", utils.ErrEmailAlreadyInUse
	}

	// Não regenerar o código nem estender a expiração; apenas reenviar o existente
	return destEmail, code, nil
}
