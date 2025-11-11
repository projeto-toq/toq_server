package userservices

import (
	"context"
	"database/sql"
	"strings"

	"errors"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
	"golang.org/x/crypto/bcrypt"
)

// SignIn autentica um usuário e retorna tokens de acesso
func (us *userService) SignIn(ctx context.Context, nationalID string, password string, deviceToken string, deviceID string) (tokens usermodel.Tokens, err error) {
	return us.SignInWithContext(ctx, nationalID, password, deviceToken, deviceID, "", "")
}

// SignInWithContext autentica um usuário com contexto de requisição completo
func (us *userService) SignInWithContext(ctx context.Context, nationalID string, password string, deviceToken string, deviceID string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Basic input validation
	if strings.TrimSpace(nationalID) == "" || strings.TrimSpace(password) == "" {
		logger.Warn("auth.signin.bad_request", "has_national_id", strings.TrimSpace(nationalID) != "", "has_password", strings.TrimSpace(password) != "")
		err = utils.BadRequest("nationalID and password are required")
		return
	}

	ctx, trimmedToken, trimmedDeviceID, derr := us.sanitizeDeviceContext(ctx, deviceToken, deviceID, "auth.signin")
	if derr != nil {
		err = derr
		return
	}

	// Normalize nationalID to digits-only (CPF/CNPJ) for consistent lookup
	nationalID = validators.OnlyDigits(nationalID)

	// Debug: valores recebidos
	if did, _ := ctx.Value(globalmodel.DeviceIDKey).(string); true {
		logger.Debug("auth.signin_with_context.debug",
			// Avoid logging PII values such as national_id
			"has_national_id", nationalID != "",
			"has_device_token", trimmedToken != "",
			"ip", ipAddress,
			"user_agent", userAgent,
			"ctx_device_id", did,
		)
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("auth.signin.tx_start_error", "error", txErr)
		return tokens, utils.InternalError("Failed to start transaction")
	}

	// Track transaction status to prevent double rollback
	txCommitted := false
	defer func() {
		if err != nil && !txCommitted {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("auth.signin.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	tokens, err = us.signIn(ctx, tx, nationalID, password, trimmedToken, trimmedDeviceID, ipAddress, userAgent)

	// CRITICAL: If authentication failed, commit the transaction to persist failed attempt tracking
	// The signIn function already processed and logged the failed attempt
	if err != nil {
		// Check if it's a domain error that should persist tracking data
		if domainErr, ok := err.(utils.DomainError); ok {
			code := domainErr.Code()
			// Commit for authentication (401), authorization (403), or user blocked (423) errors
			if code == 401 || code == 403 || code == 423 {
				if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
					utils.SetSpanError(ctx, commitErr)
					logger.Error("auth.signin.tx_commit_after_auth_error", "error", commitErr)
					// Return internal error instead of auth error since commit failed
					return tokens, utils.InternalError("Failed to process authentication")
				}
				txCommitted = true
			}
		}
		return
	}

	// Commit the transaction on successful authentication
	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("auth.signin.tx_commit_error", "error", commitErr)
		err = utils.InternalError("Failed to commit transaction")
		return
	}
	txCommitted = true

	return
}

func (us *userService) signIn(ctx context.Context, tx *sql.Tx, nationalID string, password string, deviceToken string, deviceID string, _ string, _ string) (tokens usermodel.Tokens, err error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Repository now returns user WITH active role in single query
	user, err := us.repo.GetUserByNationalID(ctx, tx, nationalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Log da tentativa com nationalID inexistente (sem expor se existe ou não)
			logger.Debug("auth.signin.invalid_credentials", "security", true)
			err = utils.AuthenticationError("Invalid credentials")
			return
		}
		utils.SetSpanError(ctx, err)
		logger.Error("auth.signin.get_user_by_nid_error", "national_id", nationalID, "error", err)
		err = utils.InternalError("Failed to validate credentials")
		return
	}

	userID := user.GetID()

	// Verificação única de bloqueio temporário ANTES de qualquer validação
	isBlocked, err := us.IsUserTempBlockedWithTx(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("auth.signin.check_temp_block_error", "user_id", userID, "error", err)
		err = utils.InternalError("Failed to validate user status")
		return
	}

	if isBlocked {
		// SECURITY: Return generic error to prevent account enumeration
		// User will receive notification via email/SMS about the block
		logger.Warn("auth.signin.user_blocked", "security", true, "user_id", userID)
		err = utils.AuthenticationError("Invalid credentials")
		return
	}

	// Validate domain invariant: repository already populated active role
	activeRole := user.GetActiveRole()
	if activeRole == nil {
		// Log da tentativa sem role ativa
		logger.Warn("auth.signin.no_active_role", "security", true, "user_id", userID)
		err = utils.AuthorizationError("No active user roles")
		return
	}

	// Comparar a senha fornecida com o hash armazenado (bcrypt)
	if bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(password)) != nil {
		// Process failed attempt and persist tracking record
		err = us.processFailedSigninAttempt(ctx, tx, userID)
		if err != nil {
			return
		}

		// Log da tentativa com credenciais inválidas
		// Transaction will be committed by SignInWithContext to persist failed attempt counter
		logger.Debug("auth.signin.invalid_credentials", "security", true, "user_id", userID)
		err = utils.AuthenticationError("Invalid credentials")
		return
	}

	// Limpa registros de tentativas erradas em caso de sucesso
	_, err = us.repo.DeleteWrongSignInByUserID(ctx, tx, userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.SetSpanError(ctx, err)
			logger.Error("auth.signin.delete_wrong_signin_error", "user_id", userID, "error", err)
			err = utils.InternalError("Failed to clear signin attempts")
			return
		}
	}

	// Remove bloqueio temporário se login foi bem-sucedido
	err = us.clearTemporaryBlockOnSuccess(ctx, tx, userID)
	if err != nil {
		// Não falha o login por problema de desbloqueio, apenas loga
		logger.Warn("auth.signin.clear_temp_block_failed", "user_id", userID, "error", err)
	}

	// Adiciona device token vinculado ao deviceID sanitizado
	if deviceToken != "" {
		logger.Debug("auth.signin.device_token.add_for_device", "device_id", deviceID, "user_id", userID)
		if _, errAdd := us.deviceTokenRepo.AddTokenForDevice(userID, deviceID, deviceToken, nil); errAdd != nil {
			logger.Warn("auth.signin.device_token_add_for_device_failed", "user_id", userID, "device_id", deviceID, "error", errAdd)
		}
	}

	// Gera os tokens
	tokens, err = us.CreateTokens(ctx, tx, user, false)
	if err != nil {
		return
	}

	// Log do sucesso no signin
	logger.Info("auth.signin.success", "security", true, "user_id", userID)

	// Note: Last activity é rastreada automaticamente pelo AuthInterceptor → Redis → Batch worker
	// Não é necessário chamar UpdateUserLastActivity diretamente

	return
}

// clearTemporaryBlockOnSuccess remove bloqueio temporário após login bem-sucedido
func (us *userService) clearTemporaryBlockOnSuccess(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	isBlocked, err := us.IsUserTempBlockedWithTx(ctx, tx, userID)
	if err != nil {
		return err
	}

	if isBlocked {
		err = us.repo.UnblockUser(ctx, tx, userID)
		if err != nil {
			return err
		}

		// Log do desbloqueio
		logger.Info("auth.signin.user_unblocked", "security", true, "user_id", userID)
	}

	return nil
}
