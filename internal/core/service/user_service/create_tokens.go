package userservices

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/giulio-alfieri/toq_server/internal/core/events"
	sessionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/session_model"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/google/uuid"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) CreateTokens(ctx context.Context, tx *sql.Tx, user usermodel.UserInterface, expired bool) (tokens usermodel.Tokens, err error) {
	// Método interno: não iniciar novo tracer; reutilizar ctx
	ctx = utils.ContextWithLogger(ctx)

	tokens = usermodel.Tokens{}

	if tx == nil {
		utils.LoggerFromContext(ctx).Error("user.create_tokens.tx_nil")
		return tokens, utils.InternalError("Transaction is nil")
	}

	secret := globalmodel.GetJWTSecret()

	expires := us.GetTokenExpiration(expired)

	accessToken, err := us.CreateAccessToken(secret, user, expires)
	if err != nil {
		// CreateAccessToken retorna DomainError em ausência de active role; apenas propague
		return tokens, err
	}
	tokens.AccessToken = accessToken

	jti := uuid.New().String()
	err = us.CreateRefreshToken(expired, user, &tokens, jti)
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
		// Set device ID if present
		if did, ok := ctx.Value(globalmodel.DeviceIDKey).(string); ok {
			s.SetDeviceID(did)
		}
		// DeviceID placeholder: can be provided via metadata in future
		if err := us.sessionRepo.CreateSession(ctx, tx, s); err != nil {
			// Persistência de sessão é infra; logar WARN e seguir (não falha emissão de tokens)
			utils.LoggerFromContext(ctx).Warn("user.create_tokens.persist_session_failed", "err", err)
		} else {
			us.globalService.GetEventBus().Publish(events.SessionEvent{Type: events.SessionCreated, UserID: s.GetUserID(), SessionID: ptrInt64(s.GetID()), DeviceID: s.GetDeviceID()})
		}
	}

	// Note: Last activity is now tracked automatically by AuthInterceptor → Redis → Batch worker
	// No need for direct UpdateUserLastActivity call

	return
}

// ptrInt64 returns a pointer to an int64 value
func ptrInt64(v int64) *int64 { return &v }
