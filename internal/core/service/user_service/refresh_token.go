package userservices

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log/slog"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	metricRefreshSuccess = prometheus.NewCounter(prometheus.CounterOpts{Name: "auth_refresh_success_total", Help: "Total successful refresh operations"})
	metricRefreshReuse   = prometheus.NewCounter(prometheus.CounterOpts{Name: "auth_refresh_reuse_detected_total", Help: "Total detected refresh token reuse incidents"})
	metricRefreshExpired = prometheus.NewCounter(prometheus.CounterOpts{Name: "auth_refresh_absolute_expired_total", Help: "Total refresh attempts blocked due to absolute expiry"})
)

func init() {
	prometheus.MustRegister(metricRefreshSuccess, metricRefreshReuse, metricRefreshExpired)
}

func (us *userService) RefreshTokens(ctx context.Context, refresh string) (tokens usermodel.Tokens, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Start transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	tokens, err = us.refreshToken(ctx, tx, refresh)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	return
}

func (us *userService) refreshToken(ctx context.Context, tx *sql.Tx, refresh string) (tokens usermodel.Tokens, err error) {
	userID, err := validateRefreshToken(refresh)
	if err != nil {
		return
	}

	// Hash incoming refresh to locate session
	sum := sha256.Sum256([]byte(refresh))
	hash := hex.EncodeToString(sum[:])
	session, sessErr := us.sessionRepo.GetActiveSessionByRefreshHash(ctx, hash)
	if sessErr != nil {
		// If session not found treat as invalid refresh
		slog.Warn("auth.refresh.invalid_session", "error", sessErr, "refresh_hash_prefix", hash[:8])
		return tokens, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	// Optional: detect reuse (rotated_at set means token already used)
	if session.GetRotatedAt() != nil {
		_ = us.sessionRepo.RevokeSessionsByUserID(ctx, session.GetUserID())
		metricRefreshReuse.Inc()
		slog.Warn("auth.refresh.reuse_detected", "user_id", session.GetUserID(), "session_id", session.GetID())
		return tokens, status.Error(codes.Unauthenticated, "refresh token reuse detected")
	}

	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}
	if user.GetActiveRole().GetStatus() == usermodel.StatusBlocked {
		err = status.Errorf(codes.PermissionDenied, "User is blocked")
		return
	}

	// Enforce absolute expiry if set
	if !session.GetAbsoluteExpiresAt().IsZero() && time.Now().UTC().After(session.GetAbsoluteExpiresAt()) {
		_ = us.sessionRepo.RevokeSessionsByUserID(ctx, session.GetUserID())
		metricRefreshExpired.Inc()
		slog.Info("auth.refresh.absolute_expired", "user_id", session.GetUserID(), "session_id", session.GetID())
		return tokens, status.Error(codes.Unauthenticated, "session expired")
	}

	// Pass chain info through context so CreateTokens can continue absolute expiry and increment rotation counter
	ctx = context.WithValue(ctx, globalmodel.SessionAbsoluteExpiryKey, session.GetAbsoluteExpiresAt())
	ctx = context.WithValue(ctx, globalmodel.SessionRotationCounterKey, session.GetRotationCounter())

	// Enforce max rotations
	if session.GetRotationCounter() >= globalmodel.GetMaxSessionRotations() {
		_ = us.sessionRepo.RevokeSessionsByUserID(ctx, session.GetUserID())
		slog.Warn("auth.refresh.rotation_limit_exceeded", "user_id", session.GetUserID(), "session_id", session.GetID(), "rotation_counter", session.GetRotationCounter())
		return tokens, status.Error(codes.Unauthenticated, "session rotation limit reached")
	}

	// Issue new tokens (new session persisted with incremented rotation counter)
	tokens, err = us.CreateTokens(ctx, tx, user, false)
	if err != nil {
		return
	}

	// Mark old session rotated
	if err = us.sessionRepo.UpdateSessionRotation(ctx, session.GetID(), session.GetRotationCounter(), time.Now().UTC()); err != nil {
		slog.Warn("auth.refresh.update_rotation_failed", "session_id", session.GetID(), "err", err)
	} else {
		slog.Info("auth.refresh.ok", "user_id", session.GetUserID(), "prev_session_id", session.GetID())
	}
	metricRefreshSuccess.Inc()

	return
}
