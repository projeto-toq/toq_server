package userservices

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/giulio-alfieri/toq_server/internal/core/events"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

var (
	metricRefreshSuccess = prometheus.NewCounter(prometheus.CounterOpts{Name: "auth_refresh_success_total", Help: "Total successful refresh operations"})
	metricRefreshReuse   = prometheus.NewCounter(prometheus.CounterOpts{Name: "auth_refresh_reuse_detected_total", Help: "Total detected refresh token reuse incidents"})
	metricRefreshExpired = prometheus.NewCounter(prometheus.CounterOpts{Name: "auth_refresh_absolute_expired_total", Help: "Total refresh attempts blocked due to absolute expiry"})
	metricSessionRotated = prometheus.NewCounter(prometheus.CounterOpts{Name: "auth_session_rotated_total", Help: "Total number of session rotations"})
)

func init() {
	prometheus.MustRegister(metricRefreshSuccess, metricRefreshReuse, metricRefreshExpired, metricSessionRotated)
}

func (us *userService) RefreshTokens(ctx context.Context, refresh string) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// Start transaction
	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("auth.refresh.tx_start_error", "error", txErr)
		return tokens, utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("auth.refresh.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	tokens, err = us.refreshToken(ctx, tx, refresh)
	if err != nil {
		return
	}

	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("auth.refresh.tx_commit_error", "error", commitErr)
		err = utils.InternalError("Failed to commit transaction")
		return
	}

	return
}

func (us *userService) refreshToken(ctx context.Context, tx *sql.Tx, refresh string) (tokens usermodel.Tokens, err error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	userID, err := validateRefreshToken(refresh)
	if err != nil {
		// Erro 401 de token inválido
		return tokens, err
	}

	// Hash incoming refresh to locate session
	sum := sha256.Sum256([]byte(refresh))
	hash := hex.EncodeToString(sum[:])
	session, sessErr := us.sessionRepo.GetActiveSessionByRefreshHash(ctx, tx, hash)
	if sessErr != nil {
		// Sessão não encontrada para o refresh apresentado
		logger.Warn("auth.refresh.invalid_session", "error", sessErr, "refresh_hash_prefix", hash[:8])
		// Captura origem ao criar o erro de domínio
		return tokens, utils.WrapDomainErrorWithSource(utils.ErrRefreshSessionNotFound)
	}

	// Optional: detect reuse (rotated_at set means token already used)
	if session.GetRotatedAt() != nil {
		_ = us.sessionRepo.RevokeSessionsByUserID(ctx, tx, session.GetUserID())
		us.globalService.GetEventBus().Publish(events.SessionEvent{Type: events.SessionsRevoked, UserID: session.GetUserID(), DeviceID: session.GetDeviceID()})
		metricRefreshReuse.Inc()
		logger.Warn("auth.refresh.reuse_detected", "user_id", session.GetUserID(), "session_id", session.GetID())
		return tokens, utils.WrapDomainErrorWithSource(utils.ErrRefreshTokenReuseDetected)
	}

	// Carrega usuário e tenta obter role ativa (necessária para access token)
	// Carrega usuário com active role via Service (invariável: requer active role)
	user, err := us.GetUserByIDWithTx(ctx, tx, userID)
	if err != nil {
		return
	}

	// Segurança adicional: o userID do JWT deve bater com a sessão carregada
	if session.GetUserID() != userID {
		_ = us.sessionRepo.RevokeSessionsByUserID(ctx, tx, session.GetUserID())
		us.globalService.GetEventBus().Publish(events.SessionEvent{Type: events.SessionsRevoked, UserID: session.GetUserID(), DeviceID: session.GetDeviceID()})
		logger.Warn("auth.refresh.user_mismatch", "jwt_user_id", userID, "session_user_id", session.GetUserID(), "session_id", session.GetID())
		return tokens, utils.WrapDomainErrorWithSource(utils.ErrInvalidRefreshToken)
	}

	// Enforce absolute expiry if set
	if !session.GetAbsoluteExpiresAt().IsZero() && time.Now().UTC().After(session.GetAbsoluteExpiresAt()) {
		_ = us.sessionRepo.RevokeSessionsByUserID(ctx, tx, session.GetUserID())
		us.globalService.GetEventBus().Publish(events.SessionEvent{Type: events.SessionsRevoked, UserID: session.GetUserID(), DeviceID: session.GetDeviceID()})
		metricRefreshExpired.Inc()
		logger.Info("auth.refresh.absolute_expired", "user_id", session.GetUserID(), "session_id", session.GetID())
		return tokens, utils.WrapDomainErrorWithSource(utils.ErrRefreshTokenExpired)
	}

	// Pass chain info through context so CreateTokens can continue absolute expiry and increment rotation counter
	ctx = context.WithValue(ctx, globalmodel.SessionAbsoluteExpiryKey, session.GetAbsoluteExpiresAt())
	ctx = context.WithValue(ctx, globalmodel.SessionRotationCounterKey, session.GetRotationCounter())

	// Enforce max rotations
	if session.GetRotationCounter() >= globalmodel.GetMaxSessionRotations() {
		_ = us.sessionRepo.RevokeSessionsByUserID(ctx, tx, session.GetUserID())
		us.globalService.GetEventBus().Publish(events.SessionEvent{Type: events.SessionsRevoked, UserID: session.GetUserID(), DeviceID: session.GetDeviceID()})
		logger.Warn("auth.refresh.rotation_limit_exceeded", "user_id", session.GetUserID(), "session_id", session.GetID(), "rotation_counter", session.GetRotationCounter())
		return tokens, utils.WrapDomainErrorWithSource(utils.ErrRefreshRotationLimitExceeded)
	}

	// Issue new tokens (new session persisted with incremented rotation counter)
	// Enrich trace context com atributos úteis
	// Comentário em português: atributos no contexto ajudam a rastrear a cadeia de rotação
	tokens, err = us.CreateTokens(ctx, tx, user, false)
	if err != nil {
		// Se falhar por ausência de role ativa nesta etapa, é um erro de domínio claro
		// e não deve ser reportado como refresh inválido. O handler retornará o código adequado.
		return
	}

	// Marca a sessão anterior como rotacionada (rotated_at), evitando reutilização do refresh token antigo
	// Comentário em português: esta marcação permite detectar "reuse" com base em rotated_at != nil
	if err = us.sessionRepo.MarkSessionRotated(ctx, tx, session.GetID()); err != nil {
		logger.Warn("auth.refresh.mark_rotated_failed", "session_id", session.GetID(), "error", err)
	}

	// Atualiza metadados da sessão antiga (contador e last_refresh_at) para fins de auditoria
	if err = us.sessionRepo.UpdateSessionRotation(ctx, tx, session.GetID(), session.GetRotationCounter(), time.Now().UTC()); err != nil {
		logger.Warn("auth.refresh.update_rotation_failed", "session_id", session.GetID(), "error", err)
	} else {
		// Publish SessionRotated for previous session
		sid := session.GetID()
		us.globalService.GetEventBus().Publish(events.SessionEvent{Type: events.SessionRotated, UserID: session.GetUserID(), SessionID: &sid, DeviceID: session.GetDeviceID()})
		metricSessionRotated.Inc()
		logger.Info("auth.refresh.ok", "user_id", session.GetUserID(), "prev_session_id", session.GetID())
	}
	// Métricas: sucesso
	metricRefreshSuccess.Inc()

	return
}
