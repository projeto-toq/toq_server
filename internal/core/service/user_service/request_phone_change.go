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

// RequestPhoneChange starts the phone change flow by generating a validation code
// and persisting the new phone as pending. If the new phone equals the current one,
// the operation is a no-op (no pending created, no notification). The user ID is
// read from context (SSOT). The phone is normalized to E.164.
func (us *userService) RequestPhoneChange(ctx context.Context, newPhone string) (err error) {
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

	// Normalize to E.164 (also validates)
	if newPhone, err = validators.NormalizeToE164(newPhone); err != nil {
		// Map validator error to a domain validation error
		return utils.ValidationError("phone", err.Error())
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		slog.Error("phone_change.request.tx_start_error", "error", txErr)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("phone_change.request.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	user, validation, err := us.requestPhoneChange(ctx, tx, userID, newPhone)
	if err != nil {
		return
	}

	// Commit the transaction BEFORE sending notification
	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		slog.Error("phone_change.request.tx_commit_error", "error", commitErr)
		return utils.InternalError("Failed to commit transaction")
	}

	// Se não houve pendência criada (mesmo telefone do atual), retornar sucesso sem notificar
	if validation == nil || validation.GetPhoneCode() == "" {
		return nil
	}

	// Usar o sistema unificado de notificação
	notificationService := us.globalService.GetUnifiedNotificationService()
	smsRequest := globalservice.NotificationRequest{
		Type: globalservice.NotificationTypeSMS,
		To:   validation.GetNewPhone(),
		Body: "TOQ - Seu código de validação: " + validation.GetPhoneCode(),
	}

	notifyErr := notificationService.SendNotification(ctx, smsRequest)
	if notifyErr != nil {
		// Log sem afetar operação principal (commit já feito)
		utils.SetSpanError(ctx, notifyErr)
		slog.Error("phone_change.request.notification_error", "user_id", user.GetID(), "error", notifyErr)
	}

	return
}

func (us *userService) requestPhoneChange(ctx context.Context, tx *sql.Tx, id int64, phone string) (user usermodel.UserInterface, validation usermodel.ValidationInterface, err error) {

	user, err = us.repo.GetUserByID(ctx, tx, id)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("phone_change.request.read_user_error", "error", err, "user_id", id)
		return
	}

	// No-op: se o novo telefone for igual ao atual, não criar pendência
	if user.GetPhoneNumber() == phone {
		return user, nil, nil
	}

	// If phone already in use by another user
	if exist, verr := us.repo.ExistsPhoneForAnotherUser(ctx, tx, phone, user.GetID()); verr != nil {
		utils.SetSpanError(ctx, verr)
		slog.Error("phone_change.request.exists_phone_error", "error", verr, "user_id", user.GetID())
		return nil, nil, verr
	} else if exist {
		return nil, nil, utils.ErrPhoneAlreadyInUse
	}

	// set the user validation as pending for phone
	validation, err = us.repo.GetUserValidations(ctx, tx, user.GetID())
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.SetSpanError(ctx, err)
			slog.Error("phone_change.request.read_validations_error", "error", err, "user_id", user.GetID())
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
		utils.SetSpanError(ctx, err)
		slog.Error("phone_change.request.update_validations_error", "error", err, "user_id", user.GetID())
		return
	}

	// Note: SendNotification moved to after transaction commit
	// Note: Last activity is now tracked automatically by AuthInterceptor → Redis → Batch worker

	return
}
