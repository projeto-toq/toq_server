package userservices

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/projeto-toq/toq_server/internal/core/events"
	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	auditservice "github.com/projeto-toq/toq_server/internal/core/service/audit_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	metricSignoutTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_signout_total",
		Help: "Total signout requests by mode",
	}, []string{"mode"})
	metricSignoutSessionsRevoked = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_signout_sessions_revoked_total",
		Help: "Total sessions revoked during signout by mode",
	}, []string{"mode"})
	metricSignoutDeviceTokensRemoved = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_signout_device_tokens_removed_total",
		Help: "Device tokens removed during signout by mode/method/result",
	}, []string{"mode", "method", "result"})
)

func init() {
	prometheus.MustRegister(metricSignoutTotal, metricSignoutSessionsRevoked, metricSignoutDeviceTokensRemoved)
}

func (us *userService) SignOut(ctx context.Context, deviceToken, refreshToken, deviceID string) (err error) {
	// Obter o ID do usuário do contexto (SSOT)
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return utils.AuthenticationError("")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	trimmedDeviceToken := strings.TrimSpace(deviceToken)
	trimmedRefreshToken := strings.TrimSpace(refreshToken)
	trimmedDeviceID := strings.TrimSpace(deviceID)
	tokenJTI, _ := ctx.Value(globalmodel.AccessTokenJTIKey).(string)
	accessExp, _ := ctx.Value(globalmodel.AccessTokenExpiresAtKey).(time.Time)

	if trimmedDeviceID != "" {
		if _, parseErr := uuid.Parse(trimmedDeviceID); parseErr != nil {
			logger.Warn("auth.signout.invalid_device_id", "device_id", trimmedDeviceID)
			return utils.BadRequest("deviceID must be a valid UUID")
		}
		// Propagar device ID sanitizado para o contexto
		ctx = context.WithValue(ctx, globalmodel.DeviceIDKey, trimmedDeviceID)
	}

	tx, txErr := us.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("auth.signout.tx_start_error", "error", txErr)
		return utils.InternalError("Failed to start transaction")
	}
	defer func() {
		if err != nil {
			if rbErr := us.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("auth.signout.tx_rollback_error", "error", rbErr)
			}
		}
	}()

	err = us.signOut(ctx, tx, userID, trimmedDeviceToken, trimmedRefreshToken, trimmedDeviceID)
	if err != nil {
		return
	}

	if tokenJTI != "" && !accessExp.IsZero() && us.tokenBlocklist != nil {
		now := time.Now().UTC()
		ttlSeconds := int64(accessExp.Sub(now).Seconds())
		if ttlSeconds < 0 {
			ttlSeconds = 0
		}
		if errBlk := us.tokenBlocklist.Add(ctx, tokenJTI, ttlSeconds); errBlk != nil {
			utils.SetSpanError(ctx, errBlk)
			logger.Warn("auth.signout.blocklist_failed", "error", errBlk)
			return utils.InternalError("Failed to revoke access token")
		}
	}

	// Commit the transaction
	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		utils.SetSpanError(ctx, commitErr)
		logger.Error("auth.signout.tx_commit_error", "error", commitErr)
		return utils.InternalError("Failed to commit transaction")
	}
	return
}

func (us *userService) signOut(ctx context.Context, tx *sql.Tx, userID int64, deviceToken, refreshToken, deviceID string) (err error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	auditTarget := auditmodel.AuditTarget{Type: auditmodel.TargetUser, ID: userID}
	sessionID := int64(0)
	// Determine single-session vs global logout
	single := refreshToken != "" || deviceToken != "" || deviceID != ""
	mode := "global"
	if single {
		mode = "single"
	}
	metricSignoutTotal.WithLabelValues(mode).Inc()

	targetDeviceID := deviceID

	if single {
		// Revoke only the session matching the refresh token (if provided & valid)
		if refreshToken != "" && us.sessionRepo != nil {
			// Hash and find session
			sum := sha256.Sum256([]byte(refreshToken))
			hash := hex.EncodeToString(sum[:])
			// Best-effort fetch active session
			if session, sessErr := us.sessionRepo.GetActiveSessionByRefreshHash(ctx, tx, hash); sessErr == nil {
				if revokeErr := us.sessionRepo.RevokeSession(ctx, tx, session.GetID()); revokeErr != nil {
					logger.Warn("auth.signout.single.revoke_failed", "user_id", userID, "session_id", session.GetID(), "error", revokeErr)
				}
				// Publish SessionsRevoked (single)
				us.globalService.GetEventBus().Publish(events.SessionEvent{Type: events.SessionsRevoked, UserID: userID, DeviceID: session.GetDeviceID()})
				// Count revoked session (best-effort; consider success when no revoke error)
				metricSignoutSessionsRevoked.WithLabelValues(mode).Inc()

				if session.GetDeviceID() != "" {
					targetDeviceID = session.GetDeviceID()
				}
				sessionID = session.GetID()
				auditTarget = auditmodel.AuditTarget{Type: auditmodel.TargetSession, ID: sessionID}
			}
		}
		// Remove device tokens
		if deviceToken != "" {
			// Prefer explicit token removal if provided
			if errTok := us.repo.RemoveDeviceToken(ctx, tx, userID, deviceToken); errTok != nil {
				logger.Warn("auth.signout.single.device_token_delete_failed", "user_id", userID, "error", errTok)
				metricSignoutDeviceTokensRemoved.WithLabelValues(mode, "token", "error").Inc()
			} else {
				metricSignoutDeviceTokensRemoved.WithLabelValues(mode, "token", "success").Inc()
			}
		}
		if targetDeviceID != "" {
			if errTok := us.repo.RemoveDeviceTokensByDeviceID(ctx, tx, userID, targetDeviceID); errTok != nil {
				logger.Warn("auth.signout.single.device_tokens_by_device_delete_failed", "user_id", userID, "device_id", targetDeviceID, "error", errTok)
				metricSignoutDeviceTokensRemoved.WithLabelValues(mode, "device", "error").Inc()
			} else {
				metricSignoutDeviceTokensRemoved.WithLabelValues(mode, "device", "success").Inc()
			}
		}
		// Audit (best-effort, não bloqueia signout)
		auditRecord := auditservice.BuildRecordFromContext(
			ctx,
			userID,
			auditTarget,
			auditmodel.OperationAuthSignout,
			map[string]any{
				"mode":                    mode,
				"device_id":               targetDeviceID,
				"refresh_token_present":   refreshToken != "",
				"device_token_present":    deviceToken != "",
				"session_revocation_mode": "single",
				"session_id":              sessionID,
			},
		)
		if errAudit := us.auditService.RecordChange(ctx, tx, auditRecord); errAudit != nil {
			logger.Warn("auth.signout.audit_error", "error", errAudit, "mode", mode)
		}
	} else {
		// Global: revoke all sessions
		if us.sessionRepo != nil {
			if revokeAllErr := us.sessionRepo.RevokeSessionsByUserID(ctx, tx, userID); revokeAllErr != nil {
				logger.Warn("auth.signout.global.revoke_failed", "user_id", userID, "error", revokeAllErr)
			} else {
				metricSignoutSessionsRevoked.WithLabelValues(mode).Inc()
			}
			// Publish SessionsRevoked (global)
			us.globalService.GetEventBus().Publish(events.SessionEvent{Type: events.SessionsRevoked, UserID: userID})
		}
		// Remove all device tokens via repository
		if errAll := us.repo.RemoveAllDeviceTokensByUserID(ctx, tx, userID); errAll != nil {
			logger.Warn("auth.signout.global.device_tokens_delete_failed", "user_id", userID, "error", errAll)
			metricSignoutDeviceTokensRemoved.WithLabelValues(mode, "all", "error").Inc()
		} else {
			metricSignoutDeviceTokensRemoved.WithLabelValues(mode, "all", "success").Inc()
		}
		auditRecord := auditservice.BuildRecordFromContext(
			ctx,
			userID,
			auditmodel.AuditTarget{Type: auditmodel.TargetUser, ID: userID},
			auditmodel.OperationAuthSignout,
			map[string]any{
				"mode":                    mode,
				"session_revocation_mode": "global",
				"device_tokens_removed":   true,
			},
		)
		if errAudit := us.auditService.RecordChange(ctx, tx, auditRecord); errAudit != nil {
			logger.Warn("auth.signout.audit_error", "error", errAudit, "mode", mode)
		}
	}
	return
}

// Updated logout strategy: single (refresh/device) => revoke only that session & token; else global.
