package userservices

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) SignOut(ctx context.Context, userID int64, deviceToken, refreshToken string) (err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	err = us.signOut(ctx, tx, userID, deviceToken, refreshToken)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	// Commit the transaction
	if commitErr := us.globalService.CommitTransaction(ctx, tx); commitErr != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return commitErr
	}
	return
}

func (us *userService) signOut(ctx context.Context, tx *sql.Tx, userID int64, deviceToken, refreshToken string) (err error) {
	// Determine single-session vs global logout
	single := refreshToken != "" || deviceToken != ""

	if single {
		// Revoke only the session matching the refresh token (if provided & valid)
		if refreshToken != "" && us.sessionRepo != nil {
			// Hash and find session
			sum := sha256.Sum256([]byte(refreshToken))
			hash := hex.EncodeToString(sum[:])
			// Best-effort fetch active session
			if session, sessErr := us.sessionRepo.GetActiveSessionByRefreshHash(ctx, tx, hash); sessErr == nil {
				if revokeErr := us.sessionRepo.RevokeSession(ctx, tx, session.GetID()); revokeErr != nil {
					slog.Warn("auth.signout.single.revoke_failed", "user_id", userID, "session_id", session.GetID(), "err", revokeErr)
				}
			}
		}
		// Remove only this device token via repository
		if deviceToken != "" {
			if errTok := us.repo.RemoveDeviceToken(ctx, tx, userID, deviceToken); errTok != nil {
				slog.Warn("auth.signout.single.device_token_delete_failed", "user_id", userID, "err", errTok)
			}
		}
		// Audit
		_ = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Logout (single device)")
	} else {
		// Global: revoke all sessions
		if us.sessionRepo != nil {
			if revokeAllErr := us.sessionRepo.RevokeSessionsByUserID(ctx, tx, userID); revokeAllErr != nil {
				slog.Warn("auth.signout.global.revoke_failed", "user_id", userID, "err", revokeAllErr)
			}
		}
		// Remove all device tokens via repository
		if errAll := us.repo.RemoveAllDeviceTokens(ctx, tx, userID); errAll != nil {
			slog.Warn("auth.signout.global.device_tokens_delete_failed", "user_id", userID, "err", errAll)
		}
		_ = us.globalService.CreateAudit(ctx, tx, globalmodel.TableUsers, "Logout (global)")
	}
	return
}

// Updated logout strategy: single (refresh/device) => revoke only that session & token; else global.
