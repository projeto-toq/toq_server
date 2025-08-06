package userservices

import (
	"context"
	"database/sql"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) GetCrecisToValidateByStatus(ctx context.Context, UserRoleStatus usermodel.UserRoleStatus) (realtors []usermodel.UserInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Start a database transaction
	tx, err := us.globalService.StartTransaction(ctx)
	if err != nil {
		return
	}

	realtors, err = us.getCrecisToValidateByStatus(ctx, tx, UserRoleStatus)
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

func (us *userService) getCrecisToValidateByStatus(ctx context.Context, tx *sql.Tx, UserRoleStatus usermodel.UserRoleStatus) (realtors []usermodel.UserInterface, err error) {

	// Read the realtors user with given status from the database
	realtors, err = us.repo.GetUsersByStatus(ctx, tx, UserRoleStatus, usermodel.RoleRealtor)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			slog.Error("Failed to read realtor users with status in GetCrecisToValidateByStatus", "error", err)
			return
		}
		return nil, nil
	}

	return
}
