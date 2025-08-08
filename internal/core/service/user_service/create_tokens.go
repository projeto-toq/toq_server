package userservices

import (
	"context"
	"database/sql"
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
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

	err = us.CreateRefreshToken(expired, user.GetID(), &tokens)
	if err != nil {
		return
	}

	// Note: Last activity is now tracked automatically by AuthInterceptor → Redis → Batch worker
	// No need for direct UpdateUserLastActivity call

	return
}
