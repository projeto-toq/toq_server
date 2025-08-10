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

func (ua *UserAdapter) GetUsers(ctx context.Context, tx *sql.Tx) (users []usermodel.UserInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT id, full_name, nick_name, national_id, creci_number, creci_state, creci_validity, born_at, phone_number, email, zip_code, street, number, complement, neighborhood, city, state, photo, password, opt_status, last_activity_at, deleted, last_signin_attempt FROM users WHERE deleted = 0;`

	entities, err := ua.Read(ctx, tx, query)
	if err != nil {
		slog.Error("mysqluseradapter/GetUsers: error executing Read", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if len(entities) == 0 {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	for _, entity := range entities {
		user, err1 := userconverters.UserEntityToDomain(entity)
		if err1 != nil {
			return nil, err1
		}

		role, err1 := ua.GetActiveUserRolesByUserID(ctx, tx, user.GetID())
		if err != nil {
			return nil, err1
		}

		user.SetActiveRole(role)

		users = append(users, user)

	}

	return

}
