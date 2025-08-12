package userservices

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log/slog"
	"time"

	sessionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/session_model"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) CreateTokens(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface, expired bool) (tokens usermodel.Tokens, err error) {
	_, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tokens = usermodel.Tokens{}

	if tx == nil {
		slog.Warn("Transaction is nil on generate tokens")
		return tokens, status.Errorf(codes.Internal, "Internal server error")
	}

	secret := globalmodel.GetJWTSecret()

	expires := us.GetTokenExpiration(expired)

	accessToken, err := us.CreateAccessToken(secret, user, expires)
	if err != nil {
		return
	}
	tokens.AccessToken = accessToken

	jti := uuid.New().String()
	err = us.CreateRefreshToken(expired, user.GetID(), &tokens, jti)
	if err != nil {
		return
	}

	// Session persistence (only when not forcing expired tokens)
	if !expired && tokens.RefreshToken != "" && us.sessionRepo != nil {
		sum := sha256.Sum256([]byte(tokens.RefreshToken))
		hash := hex.EncodeToString(sum[:])
		s := sessionmodel.NewSession()
		s.SetUserID(user.GetID())
		s.SetRefreshHash(hash)
		refreshTTL := globalmodel.GetRefreshTTL()
		now := time.Now().UTC()
		absoluteOverride, _ := ctx.Value(globalmodel.SessionAbsoluteExpiryKey).(time.Time)
		prevRotation, _ := ctx.Value(globalmodel.SessionRotationCounterKey).(int)
		// Relative expiry always now+refreshTTL
		s.SetExpiresAt(now.Add(refreshTTL))
		if !absoluteOverride.IsZero() {
			// Continue chain absolute expiry
			s.SetAbsoluteExpiresAt(absoluteOverride)
			// Rotation counter increments
			s.SetRotationCounter(prevRotation + 1)
		} else {
			// First session in chain: define absolute expiry equals refreshTTL window (could be env-configured future)
			abs := now.Add(refreshTTL)
			s.SetAbsoluteExpiresAt(abs)
			s.SetRotationCounter(0)
		}
		s.SetTokenJTI(jti)
		// Enrich with context metadata if available
		if ua, ok := ctx.Value(globalmodel.UserAgentKey).(string); ok {
			s.SetUserAgent(ua)
		}
		if ip, ok := ctx.Value(globalmodel.ClientIPKey).(string); ok {
			s.SetIP(ip)
		}
		// DeviceID placeholder: can be provided via metadata in future
		if err := us.sessionRepo.CreateSession(ctx, tx, s); err != nil {
			slog.Warn("failed to persist session", "err", err)
		}
	}

	// Note: Last activity is now tracked automatically by AuthInterceptor → Redis → Batch worker
	// No need for direct UpdateUserLastActivity call

	return
}
