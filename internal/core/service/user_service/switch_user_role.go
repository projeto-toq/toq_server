package userservices

import (
	"context"
	"database/sql"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) SwitchUserRole(ctx context.Context, userID int64, userRoleID int64) (tokens usermodel.Tokens, err error) {
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

	tokens, err = us.switchUserRole(ctx, tx, userID, userRoleID)
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

func (us *userService) switchUserRole(ctx context.Context, tx *sql.Tx, userID int64, userRoleID int64) (tokens usermodel.Tokens, err error) {

	// Get user role by ID
	userRoles, err := us.repo.GetUserRolesByUserID(ctx, tx, userID)
	if err != nil {
		return
	}

	if len(userRoles) == 1 {
		err = status.Error(codes.FailedPrecondition, "User has only one role")
		return
	}

	//verify if the new role exists
	for _, role := range userRoles {
		if role.GetID() == userRoleID {
			err = nil
			break
		}
		err = status.Error(codes.InvalidArgument, "Role not found")
	}
	if err != nil {
		return
	}

	for _, role := range userRoles {
		if role.GetID() == userRoleID {
			role.SetActive(true)
		} else {
			role.SetActive(false)
		}
		err = us.repo.UpdateUserRole(ctx, tx, role)
		if err != nil {
			return
		}
	}

	user, err := us.repo.GetUserByID(ctx, tx, userID)
	if err != nil {
		return
	}

	//generate the token
	tokens, err = us.CreateTokens(ctx, tx, user, false)
	if err != nil {
		return
	}

	return
}
