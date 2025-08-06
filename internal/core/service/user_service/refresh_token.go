package userservices

import (
	"context"
	"database/sql"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	if user.GetActiveRole().GetStatus() == usermodel.StatusBlocked {
		err = status.Errorf(codes.PermissionDenied, "User is blocked")
		return
	}

	tokens, err = us.CreateTokens(ctx, tx, user, false)
	if err != nil {
		return
	}

	return
}
