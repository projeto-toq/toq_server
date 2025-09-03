package userservices

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"

	"errors"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

func (us *userService) RequestPhoneChange(ctx context.Context, userID int64, newPhone string) (err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Normalize to E.164 (also validates)
	if newPhone, err = validators.NormalizeToE164(newPhone); err != nil {
		return err
	}

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeRequest("start_tx_error")
		}
		return
	}

	user, validation, err := us.requestPhoneChange(ctx, tx, userID, newPhone)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeRequest("domain_error")
		}
		return
	}

	// Commit the transaction BEFORE sending notification
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeRequest("commit_error")
		}
		return
	}

	// Send notification (now asynchronous by default in the notification service)
	// This allows gRPC to respond immediately without needing additional goroutines

	// Usar o novo sistema unificado de notificação
	notificationService := us.globalService.GetUnifiedNotificationService()

	// Criar requisição de SMS com código de validação
	smsRequest := globalservice.NotificationRequest{
		Type: globalservice.NotificationTypeSMS,
		To:   user.GetPhoneNumber(),
		Body: "TOQ - Seu código de validação: " + validation.GetPhoneCode(),
	}

	notifyErr := notificationService.SendNotification(ctx, smsRequest)
	if notifyErr != nil {
		// Log error but don't affect the main operation since transaction is already committed
		slog.Error("Failed to send SMS notification", "userID", user.GetID(), "error", notifyErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPhoneChangeRequest("notify_error")
		}
	}
	if mp := us.globalService.GetMetrics(); mp != nil && notifyErr == nil {
		mp.IncrementPhoneChangeRequest("success")
	}

	return
}

func (us *userService) requestPhoneChange(ctx context.Context, tx *sql.Tx, id int64, phone string) (user usermodel.UserInterface, validation usermodel.ValidationInterface, err error) {

	user, err = us.repo.GetUserByID(ctx, tx, id)
	if err != nil {
		return
	}

	// Basic checks
	// if the phone is the same as current, return conflict
	if user.GetPhoneNumber() == phone {
		return nil, nil, utils.ErrSamePhoneAsCurrent
	}

	// If phone already in use by another user
	if exist, verr := us.repo.ExistsPhoneForAnotherUser(ctx, tx, phone, user.GetID()); verr != nil {
		return nil, nil, verr
	} else if exist {
		return nil, nil, utils.ErrPhoneAlreadyInUse
	}

	// set the user validation as pending for phone
	validation, err = us.repo.GetUserValidations(ctx, tx, user.GetID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		validation = usermodel.NewValidation()
	}

	validation.SetUserID(user.GetID())
	validation.SetPhoneCode(us.random6Digits())
	validation.SetPhoneCodeExp(time.Now().UTC().Add(usermodel.ValidationCodeExpiration))
	validation.SetNewPhone(phone)

	err = us.repo.UpdateUserValidations(ctx, tx, validation)
	if err != nil {
		return
	}

	// Note: SendNotification moved to after transaction commit
	// Note: Last activity is now tracked automatically by AuthInterceptor → Redis → Batch worker

	return
}
