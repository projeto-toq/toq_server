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

// ResendPhoneChangeCode regenerates the phone change code and extends its expiration.
// It requires a pending phone change; after commit, sends the new code via SMS to the new phone number.
func (us *userService) ResendPhoneChangeCode(ctx context.Context, userID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeResend("start_tx_error")
		}
		return
	}

	var destPhone, newCode string
	destPhone, newCode, err = us.resendPhoneChangeCode(ctx, tx, userID)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		if mp := us.globalService.GetMetrics(); mp != nil {
			switch err {
			case utils.ErrPhoneChangeNotPending:
				mp.IncrementPhoneChangeResend("not_pending")
			default:
				mp.IncrementPhoneChangeResend("domain_error")
			}
		}
		return
	}

	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeResend("commit_error")
		}
		return
	}

	// After commit, send SMS with the new code
	notificationService := us.globalService.GetUnifiedNotificationService()
	smsRequest := globalservice.NotificationRequest{
		Type: globalservice.NotificationTypeSMS,
		To:   destPhone,
		Body: "TOQ - Seu novo código de validação: " + newCode,
	}
	if notifyErr := notificationService.SendNotification(ctx, smsRequest); notifyErr != nil {
		slog.Error("Failed to send SMS notification (resend phone change)", "userID", userID, "error", notifyErr)
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

	// Generate new code and extend expiration
	code = us.random6Digits()
	validation.SetPhoneCode(code)
	validation.SetPhoneCodeExp(time.Now().UTC().Add(usermodel.ValidationCodeExpiration))

	if err = us.repo.UpdateUserValidations(ctx, tx, validation); err != nil {
		return "", "", err
	}
	return destPhone, code, nil
}
