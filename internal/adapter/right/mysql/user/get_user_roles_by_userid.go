package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) (roles []usermodel.UserRoleInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entities, err := ua.Read(ctx, tx, "SELECT * FROM user_roles WHERE user_id = ?;", userID)
	if err != nil {
		slog.Error("mysqluseradapter/GetUserRolesByUserID: error executing Read", "error", err)
		return nil, utils.ErrInternalServer
	}

	if len(entities) == 0 {
		return nil, utils.ErrInternalServer
	}

	for _, entity := range entities {

		role, err1 := userconverters.UserRoleEntityToDomain(entity)
		if err1 != nil {
			return nil, err1
		}

		roles = append(roles, role)
	}

	return
}
