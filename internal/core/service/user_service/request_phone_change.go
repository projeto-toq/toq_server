package userservices

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	globalservice "github.com/giulio-alfieri/toq_server/internal/core/service/global_service"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) RequestPhoneChange(ctx context.Context, userID int64, newPhone string) (err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	user, validation, err := us.requestPhoneChange(ctx, tx, userID, newPhone)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	// Commit the transaction BEFORE sending notification
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
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
	}

	return
}

func (us *userService) requestPhoneChange(ctx context.Context, tx *sql.Tx, id int64, phoone string) (user usermodel.UserInterface, validation usermodel.ValidationInterface, err error) {

	user, err = us.repo.GetUserByID(ctx, tx, id)
	if err != nil {
		return
	}

	//set the user validation as pending for phone
	validation, err = us.repo.GetUserValidations(ctx, tx, user.GetID())
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return
		}
		validation = usermodel.NewValidation()
	}

	validation.SetUserID(user.GetID())
	validation.SetPhoneCode(us.random6Digits())
	validation.SetPhoneCodeExp(time.Now().UTC().Add(usermodel.ValidationCodeExpiration))
	validation.SetNewPhone(phoone)

	err = us.repo.UpdateUserValidations(ctx, tx, validation)
	if err != nil {
		return
	}

	// Note: SendNotification moved to after transaction commit
	// Note: Last activity is now tracked automatically by AuthInterceptor → Redis → Batch worker

	return
}
