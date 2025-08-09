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

func (us *userService) RequestEmailChange(ctx context.Context, userID int64, newEmail string) (err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	user, validation, err := us.requestEmailChange(ctx, tx, userID, newEmail)
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

	// Criar requisição de email com código de validação
	emailRequest := globalservice.NotificationRequest{
		Type:    globalservice.NotificationTypeEmail,
		To:      user.GetEmail(),
		Subject: "TOQ - Confirmação de Alteração de Email",
		Body:    "Seu código de validação para alteração de email é: " + validation.GetEmailCode(),
	}

	notifyErr := notificationService.SendNotification(ctx, emailRequest)
	if notifyErr != nil {
		// Log error but don't affect the main operation since transaction is already committed
		slog.Error("Failed to send email notification", "userID", user.GetID(), "error", notifyErr)
	}

	return
}

func (us *userService) requestEmailChange(ctx context.Context, tx *sql.Tx, id int64, email string) (user usermodel.UserInterface, validation usermodel.ValidationInterface, err error) {

	user, err = us.repo.GetUserByID(ctx, tx, id)
	if err != nil {
		return
	}

	//set the user validation as pending for email
	validation, err = us.repo.GetUserValidations(ctx, tx, user.GetID())
	if err != nil {
		if status.Code(err) != codes.NotFound {
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
