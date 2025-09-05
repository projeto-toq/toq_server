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
		utils.SetSpanError(ctx, txErr)
		slog.Error("auth.signin.tx_start_error", "error", txErr)
		return tokens, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				slog.Error("auth.signin.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	tokens, err = us.signIn(ctx, tx, nationalID, password, deviceToken, ipAddress, userAgent)
	if err != nil {
		return
	}

	// Commit the transaction
	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		slog.Error("auth.signin.tx_commit_error", "error", commitErr)
		err = utils.InternalError("Failed to commit transaction")
		return
	}

	return
}

func (us *userService) signIn(ctx context.Context, tx *sql.Tx, nationalID string, password string, deviceToken string, _ string, _ string) (tokens usermodel.Tokens, err error) {
	user, err := us.repo.GetUserByNationalID(ctx, tx, nationalID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Log da tentativa com nationalID inexistente (sem expor se existe ou não)
			slog.Debug("auth.signin.invalid_credentials", "security", true)
			err = utils.AuthenticationError("Invalid credentials")
			return
		}
		utils.SetSpanError(ctx, err)
		slog.Error("auth.signin.get_user_by_nid_error", "national_id", nationalID, "error", err)
		err = utils.InternalError("Failed to validate credentials")
		return
	}

	userID := user.GetID()

	// Verificação única de bloqueio temporário ANTES de qualquer validação
	isBlocked, err := us.permissionService.IsUserTempBlockedWithTx(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("auth.signin.check_temp_block_error", "user_id", userID, "error", err)
		err = utils.InternalError("Failed to validate user status")
		return
	}

	if isBlocked {
		// Log do bloqueio
		slog.Warn("auth.signin.user_blocked", "security", true, "user_id", userID)
		err = utils.UserBlockedError("Account temporarily blocked due to too many failed attempts")
		return
	}

	// Busca a role ativa via Permission Service (há apenas uma ativa por vez)
	activeRole, err := us.permissionService.GetActiveUserRoleWithTx(ctx, tx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("auth.signin.get_active_role_error", "user_id", userID, "error", err)
		err = utils.InternalError("Failed to validate user permissions")
		return
	}
	if activeRole == nil {
		// Log da tentativa sem role ativa
		slog.Warn("auth.signin.no_active_role", "security", true, "user_id", userID)
		err = utils.AuthorizationError("No active user roles")
		return
	}
	user.SetActiveRole(activeRole)

	// Comparar a senha fornecida com o hash armazenado (bcrypt)
	if bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(password)) != nil {
		err = us.processWrongSignin(ctx, tx, user)
		if err != nil {
			return
		}

		// Log da tentativa com credenciais inválidas
		slog.Debug("auth.signin.invalid_credentials", "security", true, "user_id", userID)
		err = utils.AuthenticationError("Invalid credentials")
		return
	}

	// Limpa registros de tentativas erradas em caso de sucesso
	_, err = us.repo.DeleteWrongSignInByUserID(ctx, tx, userID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			utils.SetSpanError(ctx, err)
			slog.Error("auth.signin.delete_wrong_signin_error", "user_id", userID, "error", err)
			err = utils.InternalError("Failed to clear signin attempts")
			return
		}
	}

	// Remove bloqueio temporário se login foi bem-sucedido
	err = us.clearTemporaryBlockOnSuccess(ctx, tx, userID)
	if err != nil {
		// Não falha o login por problema de desbloqueio, apenas loga
		slog.Warn("auth.signin.clear_temp_block_failed", "user_id", userID, "error", err)
	}

	// Adiciona device token se fornecido (preferindo associação por deviceID quando disponível)
	if deviceToken != "" {
		if did, ok := ctx.Value(globalmodel.DeviceIDKey).(string); ok && did != "" {
			slog.Debug("auth.signin.device_token.add_for_device", "device_id", did, "user_id", userID)
			if errAdd := us.repo.AddTokenForDevice(ctx, tx, userID, did, deviceToken, nil); errAdd != nil {
				slog.Warn("auth.signin.device_token_add_for_device_failed", "user_id", userID, "device_id", did, "error", errAdd)
			}
		} else {
			slog.Debug("auth.signin.device_token.add_user_only", "user_id", userID)
			if errAdd := us.repo.AddDeviceToken(ctx, tx, userID, deviceToken, nil); errAdd != nil {
				slog.Warn("auth.signin.device_token_add_failed", "user_id", userID, "error", errAdd)
			}
		}
	}

	// Gera os tokens
	tokens, err = us.CreateTokens(ctx, tx, user, false)
	if err != nil {
		return
	}

	// Log do sucesso no signin
	slog.Info("auth.signin.success", "security", true, "user_id", userID)

	// Note: Last activity é rastreada automaticamente pelo AuthInterceptor → Redis → Batch worker
	// Não é necessário chamar UpdateUserLastActivity diretamente

	return
}

// processWrongSignin processa tentativas de signin incorretas com melhor logging
func (us *userService) processWrongSignin(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface) error {
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
		slog.Warn("auth.signin.user_blocked", "security", true, "user_id", userID, "attempts", wrongSignin.GetFailedAttempts())
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
		slog.Info("auth.signin.user_unblocked", "security", true, "user_id", userID)
	}

	return nil
}
