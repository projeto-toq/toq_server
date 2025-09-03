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

// RequestPasswordChange starts the password reset flow in a privacy-preserving way.
// It will not reveal whether the user exists. If the user is not found, it returns nil
// and records a metric with result=user_not_found without sending any notification.
func (us *userService) RequestPasswordChange(ctx context.Context, nationalID string) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPasswordChangeRequest("start_tx_error")
		}
		return
	}

	user, validation, err := us.requestPasswordChange(ctx, tx, nationalID)
	if err != nil {
		// Privacy-preserving path: do not reveal user existence
		if errors.Is(err, sql.ErrNoRows) {
			us.globalService.RollbackTransaction(ctx, tx)
			if mp := us.globalService.GetMetrics(); mp != nil {
				mp.IncrementPasswordChangeRequest("user_not_found")
			}
			return nil
		}
		us.globalService.RollbackTransaction(ctx, tx)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPasswordChangeRequest("domain_error")
		}
		return
	}

	// Commit before notify
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPasswordChangeRequest("commit_error")
		}
		return
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
		slog.Error("Failed to send password reset notification", "userID", user.GetID(), "error", notifyErr)
		if mp := us.globalService.GetMetrics(); mp != nil {
			mp.IncrementPasswordChangeRequest("notify_error")
		}
		return nil
	}

	if mp := us.globalService.GetMetrics(); mp != nil {
		mp.IncrementPasswordChangeRequest("success")
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
