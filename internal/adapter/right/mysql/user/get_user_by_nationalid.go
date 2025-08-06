package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (ua *UserAdapter) GetUserByNationalID(ctx context.Context, tx *sql.Tx, nationalID string) (user usermodel.UserInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entities, err := ua.Read(ctx, tx, "SELECT * FROM users WHERE national_id = ?", nationalID)
	if err != nil {
		slog.Error("mysqluseradapter/GetUserByNationalID: error executing Read", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if len(entities) == 0 {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	if len(entities) > 1 {
		slog.Error("mysqluseradapter/GetUserByNationalID: multiple users found with the same nationalID", "nationalID", nationalID)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	user, err = userconverters.UserEntityToDomain(entities[0])
	if err != nil {
		return
	}

	role, err := ua.GetActiveUserRolesByUserID(ctx, tx, user.GetID())
	if err != nil {
		return
	}

	user.SetActiveRole(role)

	return

}
