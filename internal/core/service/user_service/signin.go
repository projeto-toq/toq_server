package userservices

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
	"time"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"errors"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"golang.org/x/crypto/bcrypt"
)

// SignIn autentica um usuário e retorna tokens de acesso
func (us *userService) SignIn(ctx context.Context, nationalID string, password string, deviceToken string) (tokens usermodel.Tokens, err error) {
	return us.SignInWithContext(ctx, nationalID, password, deviceToken, "", "")
}

// SignInWithContext autentica um usuário com contexto de requisição completo
func (us *userService) SignInWithContext(ctx context.Context, nationalID string, password string, deviceToken string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Basic input validation
	if strings.TrimSpace(nationalID) == "" || strings.TrimSpace(password) == "" {
		slog.Warn("auth.signin.bad_request", "has_national_id", strings.TrimSpace(nationalID) != "", "has_password", strings.TrimSpace(password) != "")
		err = utils.BadRequest("nationalID and password are required")
		return
	}

	// Debug: valores recebidos
	if did, _ := ctx.Value(globalmodel.DeviceIDKey).(string); true {
		slog.Debug("auth.signin_with_context.debug",
			"national_id", nationalID,
			"has_device_token", deviceToken != "",
			"ip", ipAddress,
			"user_agent", userAgent,
			"ctx_device_id", did,
		)
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		slog.Error("auth.signin.tx_start_error", "err", txErr)
		return tokens, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				slog.Error("auth.signin.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	tokens, err = us.signIn(ctx, tx, nationalID, password, deviceToken, ipAddress, userAgent)
	if err != nil {
		return
	}

	// Commit the transaction
	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		slog.Error("auth.signin.tx_commit_error", "err", commitErr)
		err = utils.InternalError("Failed to commit transaction")
		return
	}

	return
}

func (us *userService) signIn(ctx context.Context, tx *sql.Tx, nationalID string, password string, deviceToken string, ipAddress string, userAgent string) (tokens usermodel.Tokens, err error) {
	user, err := us.repo.GetUserByNationalID(ctx, tx, nationalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Log da tentativa com nationalID inexistente (sem expor se existe ou não)
			us.securityLogger.LogInvalidCredentials(ctx, nationalID, ipAddress, userAgent)
			slog.Warn("Signin attempt with non-existent nationalID", "nationalID", nationalID)
			err = utils.AuthenticationError("Invalid credentials")
			return
		}
		slog.Error("Failed to get user by national ID", "nationalID", nationalID, "error", err)
		err = utils.InternalError("Failed to validate credentials")
		return
	}

	userID := user.GetID()

	// Verificação única de bloqueio temporário ANTES de qualquer validação
	isBlocked, err := us.permissionService.IsUserTempBlockedWithTx(ctx, tx, userID)
	if err != nil {
		slog.Error("Failed to check if user is temporarily blocked", "userID", userID, "error", err)
		err = utils.InternalError("Failed to validate user status")
		return
	}

	if isBlocked {
		// Log do bloqueio
		us.securityLogger.LogSigninAttempt(ctx, nationalID, &userID, false, &[]usermodel.SigninErrorType{usermodel.SigninErrorUserBlocked}[0], ipAddress, userAgent)
		slog.Warn("Signin attempt from temporarily blocked user", "userID", userID, "nationalID", nationalID)
		err = utils.UserBlockedError("Account temporarily blocked due to too many failed attempts")
		return
	}

	// Busca a role ativa via Permission Service (há apenas uma ativa por vez)
	activeRole, err := us.permissionService.GetActiveUserRoleWithTx(ctx, tx, userID)
	if err != nil {
		slog.Error("Failed to get active user role", "userID", userID, "error", err)
		err = utils.InternalError("Failed to validate user permissions")
		return
	}
	if activeRole == nil {
		// Log da tentativa sem role ativa
		us.securityLogger.LogNoActiveRoles(ctx, userID, nationalID, ipAddress, userAgent)
		slog.Warn("User has no active role", "userID", userID)
		err = utils.AuthorizationError("No active user roles")
		return
	}
	user.SetActiveRole(activeRole)

	// Comparar a senha fornecida com o hash armazenado (bcrypt)
	if bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(password)) != nil {
		err = us.processWrongSignin(ctx, tx, user, nationalID, ipAddress, userAgent)
		if err != nil {
			return
		}

		// Log da tentativa com credenciais inválidas
		us.securityLogger.LogSigninAttempt(ctx, nationalID, &userID, false, &[]usermodel.SigninErrorType{usermodel.SigninErrorInvalidCredentials}[0], ipAddress, userAgent)
		slog.Warn("Invalid password attempt", "userID", userID, "nationalID", nationalID)
		err = utils.AuthenticationError("Invalid credentials")
		return
	}

	// Limpa registros de tentativas erradas em caso de sucesso
	_, err = us.repo.DeleteWrongSignInByUserID(ctx, tx, userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.Error("Failed to delete wrong signin record", "userID", userID, "error", err)
			err = utils.InternalError("Failed to clear signin attempts")
			return
		}
	}

	// Remove bloqueio temporário se login foi bem-sucedido
	err = us.clearTemporaryBlockOnSuccess(ctx, tx, userID)
	if err != nil {
		// Não falha o login por problema de desbloqueio, apenas loga
		slog.Warn("Failed to clear temporary block after successful login", "userID", userID, "error", err)
	}

	// Adiciona device token se fornecido (preferindo associação por deviceID quando disponível)
	if deviceToken != "" {
		if did, ok := ctx.Value(globalmodel.DeviceIDKey).(string); ok && did != "" {
			slog.Debug("auth.signin.device_token.add_for_device", "device_id", did, "user_id", userID)
			if errAdd := us.repo.AddTokenForDevice(ctx, tx, userID, did, deviceToken, nil); errAdd != nil {
				slog.Warn("signIn: failed to add device token for device", "userID", userID, "deviceID", did, "err", errAdd)
			}
		} else {
			slog.Debug("auth.signin.device_token.add_user_only", "user_id", userID)
			if errAdd := us.repo.AddDeviceToken(ctx, tx, userID, deviceToken, nil); errAdd != nil {
				slog.Warn("signIn: failed to add device token", "userID", userID, "err", errAdd)
			}
		}
	}

	// Gera os tokens
	tokens, err = us.CreateTokens(ctx, tx, user, false)
	if err != nil {
		return
	}

	// Log do sucesso no signin
	us.securityLogger.LogSigninAttempt(ctx, nationalID, &userID, true, nil, ipAddress, userAgent)
	slog.Info("User signed in successfully", "userID", userID)

	// Note: Last activity é rastreada automaticamente pelo AuthInterceptor → Redis → Batch worker
	// Não é necessário chamar UpdateUserLastActivity diretamente

	return
}

// processWrongSignin processa tentativas de signin incorretas com melhor logging
func (us *userService) processWrongSignin(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface, nationalID string, ipAddress string, userAgent string) error {
	userID := user.GetID()

	wrongSignin, err := us.repo.GetWrongSigninByUserID(ctx, tx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			wrongSignin = usermodel.NewWrongSignin()
		} else {
			slog.Error("Failed to get wrong signin record", "userID", userID, "error", err)
			return utils.InternalError("Failed to check signin attempts")
		}
	}

	wrongSignin.SetUserID(userID)
	wrongSignin.SetLastAttemptAt(time.Now().UTC())
	wrongSignin.SetFailedAttempts(wrongSignin.GetFailedAttempts() + 1)

	err = us.repo.UpdateWrongSignIn(ctx, tx, wrongSignin)
	if err != nil {
		slog.Error("Failed to update wrong signin record", "userID", userID, "error", err)
		return utils.InternalError("Failed to update signin attempts")
	}

	// Verifica se deve bloquear o usuário
	if wrongSignin.GetFailedAttempts() >= usermodel.MaxWrongSigninAttempts {
		// Bloqueia usuário temporariamente usando permission service
		err = us.permissionService.BlockUserTemporarily(ctx, tx, userID, "Too many failed signin attempts")
		if err != nil {
			slog.Error("Failed to block user temporarily", "userID", userID, "error", err)
			return utils.InternalError("Failed to process security measures")
		}

		// Atualiza última tentativa de signin
		user.SetLastSignInAttempt(time.Now().UTC())
		err = us.repo.UpdateUserByID(ctx, tx, user)
		if err != nil {
			slog.Error("Failed to update user last signin attempt", "userID", userID, "error", err)
			return utils.InternalError("Failed to update user record")
		}

		// Log do bloqueio
		us.securityLogger.LogUserBlocked(ctx, userID, "Too many failed signin attempts", ipAddress, userAgent)
		slog.Warn("User temporarily blocked due to too many failed signin attempts",
			"userID", userID,
			"attempts", wrongSignin.GetFailedAttempts())
	}

	return nil
}

// clearTemporaryBlockOnSuccess remove bloqueio temporário após login bem-sucedido
func (us *userService) clearTemporaryBlockOnSuccess(ctx context.Context, tx *sql.Tx, userID int64) error {
	isBlocked, err := us.permissionService.IsUserTempBlockedWithTx(ctx, tx, userID)
	if err != nil {
		return err
	}

	if isBlocked {
		err = us.permissionService.UnblockUser(ctx, tx, userID)
		if err != nil {
			return err
		}

		// Log do desbloqueio
		us.securityLogger.LogUserUnblocked(ctx, userID, "Successful login after temporary block")
		slog.Info("User unblocked after successful login", "userID", userID)
	}

	return nil
}
