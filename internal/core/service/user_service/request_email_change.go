package userservices

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
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

	// Send notification asynchronously (non-blocking goroutine)
	// This allows gRPC to respond immediately without waiting for email delivery
	go func() {
		// Create a new context preserving Request ID but avoiding cancellation
		notifyCtx := context.Background()
		if requestID := ctx.Value(globalmodel.RequestIDKey); requestID != nil {
			notifyCtx = context.WithValue(notifyCtx, globalmodel.RequestIDKey, requestID)
		}

		notifyErr := us.globalService.SendNotification(notifyCtx, user, globalmodel.NotificationEmailChange, validation.GetEmailCode())
		if notifyErr != nil {
			// Log error but don't affect the main operation since transaction is already committed
			slog.Error("Failed to send email notification", "userID", user.GetID(), "error", notifyErr)
		}
	}()

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
