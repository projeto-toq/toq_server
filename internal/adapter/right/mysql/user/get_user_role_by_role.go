package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetUserRoleByRole(ctx context.Context, tx *sql.Tx, roleToGet usermodel.UserRole) (role usermodel.UserRoleInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM user_roles WHERE role = ?;`

	entities, err := ua.Read(ctx, tx, query, roleToGet)
	if err != nil {
		slog.Error("mysqluseradapter/GetUserRoleByRole: error executing Read", "error", err)
		return nil, utils.ErrInternalServer
	}

	if len(entities) == 0 {
		return nil, utils.ErrInternalServer
	}

	if len(entities) > 1 {
		slog.Error("mysqluseradapter/GetUserRoleByRole:  multiple roles found with the same role", "role", roleToGet)
		return nil, utils.ErrInternalServer
	}

	role, err = userconverters.UserRoleEntityToDomain(entities[0])
	if err != nil {
		return
	}

	return
}
