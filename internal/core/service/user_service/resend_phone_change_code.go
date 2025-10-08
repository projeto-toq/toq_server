package userservices

import (
	"context"
	"database/sql"
	"errors"
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

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("user.resend_phone_change_code.tx_start_error", "error", txErr)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("user.resend_phone_change_code.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	var destPhone, code string
	destPhone, code, err = us.resendPhoneChangeCode(ctx, tx, userID)
	if err != nil {
		return
	}

	if err = us.globalService.CommitTransaction(ctx, tx); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("user.resend_phone_change_code.tx_commit_error", "error", err)
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
		utils.SetSpanError(ctx, notifyErr)
		logger.Error("user.resend_phone_change_code.notification_error", "user_id", userID, "error", notifyErr)
	}
	return
}

// resendPhoneChangeCode performs the regeneration of the phone code and extends the expiration.
// Returns the destination phone (new phone) for notification purposes.
func (us *userService) resendPhoneChangeCode(ctx context.Context, tx *sql.Tx, userID int64) (destPhone string, code string, err error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	validation, err := us.repo.GetUserValidations(ctx, tx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", utils.ErrPhoneChangeNotPending
		}
		utils.SetSpanError(ctx, err)
		logger.Error("user.resend_phone_change_code.read_validations_error", "error", err, "user_id", userID)
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
		utils.SetSpanError(ctx, verr)
		logger.Error("user.resend_phone_change_code.exists_phone_error", "error", verr, "user_id", userID)
		return "", "", verr
	} else if exist {
		return "", "", utils.ErrPhoneAlreadyInUse
	}
	// Não regenerar o código nem estender a expiração; apenas reenviar o existente
	return destPhone, code, nil
}
