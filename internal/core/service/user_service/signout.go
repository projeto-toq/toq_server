package userservices

import (
	"context"
	"database/sql"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) SignOut(ctx context.Context, userID int64) (tokens usermodel.Tokens, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	tokens, err = us.signOut(ctx, tx, userID)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	// Commit the transaction
	err = us.globalService.CommitTransaction(ctx, tx)
	if err != nil {
		us.globalService.RollbackTransaction(ctx, tx)
		return
	}

	return
}

func (us *userService) signOut(ctx context.Context, tx *sql.Tx, userID int64) (tokens usermodel.Tokens, err error) {
	// Revoke all active sessions for the user
	if us.sessionRepo != nil {
		if err = us.sessionRepo.RevokeSessionsByUserID(ctx, userID); err != nil {
			slog.Warn("failed to revoke user sessions on signout", "userID", userID, "err", err)
		}
	}
	// Return empty tokens (proto may still expect structure)
	return usermodel.Tokens{}, nil
}

// Logout strategy: revoke all user's active sessions and return empty tokens.
