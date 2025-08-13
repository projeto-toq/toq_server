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

func (ua *UserAdapter) GetUserByID(ctx context.Context, tx *sql.Tx, id int64) (user usermodel.UserInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT id, full_name, nick_name, national_id, creci_number, creci_state, creci_validity, born_at, phone_number, email, zip_code, street, number, complement, neighborhood, city, state, password, opt_status, last_activity_at, deleted, last_signin_attempt FROM users WHERE id = ? AND deleted = 0;`

	entities, err := ua.Read(ctx, tx, query, id)
	if err != nil {
		slog.Error("mysqluseradapter/GetUserByID: error executing Read", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if len(entities) == 0 {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	if len(entities) > 1 {
		slog.Error("mysqluseradapter/GetUserByID: multiple users found with the same ID", "ID", id)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	user, err = userconverters.UserEntityToDomain(entities[0])
	if err != nil {
		return
	}

	role, err := ua.GetActiveUserRolesByUserID(ctx, tx, id)
	if err != nil {
		return
	}

	user.SetActiveRole(role)

	return

}
