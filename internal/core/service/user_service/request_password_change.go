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
	validators "github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

// RequestPasswordChange starts the password reset flow in a privacy-preserving way.
// It will not reveal whether the user exists. If the user is not found, it returns nil
// and records a metric with result=user_not_found without sending any notification.
func (us *userService) RequestPasswordChange(ctx context.Context, nationalID string) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	// Start transaction
	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		slog.Error("auth.request_password_change.tx_start_error", "error", txErr)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("auth.request_password_change.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	// Normalize nationalID to digits-only for consistent lookup
	nationalID = validators.OnlyDigits(nationalID)

	user, validation, err := us.requestPasswordChange(ctx, tx, nationalID)
	if err != nil {
		// Privacy-preserving path: do not reveal user existence
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return
	}

	// Commit before notify
	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		slog.Error("auth.request_password_change.tx_commit_error", "error", commitErr)
		return utils.InternalError("Failed to commit transaction")
	}

	// Send notification after commit
	notificationService := us.globalService.GetUnifiedNotificationService()
	emailRequest := globalservice.NotificationRequest{
		Type:    globalservice.NotificationTypeEmail,
		To:      user.GetEmail(),
		Subject: "TOQ - Password Reset",
		Body:    "Your password reset code is: " + validation.GetPasswordCode(),
	}

	if notifyErr := notificationService.SendNotification(ctx, emailRequest); notifyErr != nil {
		// Do not impact main operation
		utils.SetSpanError(ctx, notifyErr)
		slog.Error("auth.request_password_change.notification_error", "user_id", user.GetID(), "error", notifyErr)
		return nil
	}
	return nil
}

func (us *userService) requestPasswordChange(ctx context.Context, tx *sql.Tx, nationalID string) (user usermodel.UserInterface, validation usermodel.ValidationInterface, err error) {
	user, err = us.repo.GetUserByNationalID(ctx, tx, nationalID)
	if err != nil {
		// Propagate ErrNoRows to caller for privacy-preserving behavior
		return
	}

	// Set password validation as pending
	validation, err = us.repo.GetUserValidations(ctx, tx, user.GetID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			validation = usermodel.NewValidation()
		} else {
			return
		}
	}

	validation.SetUserID(user.GetID())
	validation.SetPasswordCode(us.random6Digits())
	validation.SetPasswordCodeExp(time.Now().UTC().Add(usermodel.ValidationCodeExpiration))

	err = us.repo.UpdateUserValidations(ctx, tx, validation)
	if err != nil {
		return
	}

	// Notification is sent after commit by caller
	return
}
