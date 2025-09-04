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

// ResendPhoneChangeCode regenerates the phone change code and extends its expiration.
// It requires a pending phone change; after commit, sends the new code via SMS to the new phone number.
func (us *userService) ResendPhoneChangeCode(ctx context.Context) (err error) {
	// Obter o ID do usuário do contexto (SSOT)
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return utils.AuthenticationError("")
	}
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		slog.Error("user.resend_phone_change_code.tx_start_error", "err", txErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeResend("start_tx_error")
		}
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("user.resend_phone_change_code.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	var destPhone, code string
	destPhone, code, err = us.resendPhoneChangeCode(ctx, tx, userID)
	if err != nil {
		if mp := us.globalService.GetMetrics(); mp != nil {
			switch err {
			case utils.ErrPhoneChangeNotPending:
				mp.IncrementPhoneChangeResend("not_pending")
			case utils.ErrPhoneChangeCodeExpired:
				mp.IncrementPhoneChangeResend("expired")
			case utils.ErrPhoneAlreadyInUse:
				mp.IncrementPhoneChangeResend("already_in_use")
			default:
				mp.IncrementPhoneChangeResend("domain_error")
			}
		}
		return
	}

	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		slog.Error("user.resend_phone_change_code.tx_commit_error", "err", err)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeResend("commit_error")
		}
		return utils.InternalError("Failed to commit transaction")
	}

	// After commit, send SMS with the new code
	notificationService := us.globalService.GetUnifiedNotificationService()
	smsRequest := globalservice.NotificationRequest{
		Type: globalservice.NotificationTypeSMS,
		To:   destPhone,
		Body: "TOQ - Seu código de validação: " + code,
	}
	if notifyErr := notificationService.SendNotification(ctx, smsRequest); notifyErr != nil {
		slog.Error("user.resend_phone_change_code.notification_error", "userID", userID, "err", notifyErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeResend("notify_error")
		}
	} else {
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeResend("success")
		}
	}
	return
}

// resendPhoneChangeCode performs the regeneration of the phone code and extends the expiration.
// Returns the destination phone (new phone) for notification purposes.
func (us *userService) resendPhoneChangeCode(ctx context.Context, tx *sql.Tx, userID int64) (destPhone string, code string, err error) {
	validation, err := us.repo.GetUserValidations(ctx, tx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", utils.ErrPhoneChangeNotPending
		}
		return "", "", err
	}

	destPhone = validation.GetNewPhone()
	if destPhone == "" {
		return "", "", utils.ErrPhoneChangeNotPending
	}
	// Deve haver um código válido ainda dentro do prazo
	code = validation.GetPhoneCode()
	if code == "" {
		return "", "", utils.ErrPhoneChangeNotPending
	}
	if validation.GetPhoneCodeExp().Before(time.Now().UTC()) {
		return "", "", utils.ErrPhoneChangeCodeExpired
	}
	// Verificar unicidade global (outros usuários não podem ter este telefone)
	if exist, verr := us.repo.ExistsPhoneForAnotherUser(ctx, tx, destPhone, userID); verr != nil {
		return "", "", verr
	} else if exist {
		return "", "", utils.ErrPhoneAlreadyInUse
	}
	// Não regenerar o código nem estender a expiração; apenas reenviar o existente
	return destPhone, code, nil
}
