package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetActiveUserRolesByUserID(ctx context.Context, tx *sql.Tx, userID int64) (role usermodel.UserRoleInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	entities, err := ua.Read(ctx, tx, "SELECT * FROM user_roles WHERE user_id = ? AND active = 1", userID)
	if err != nil {
		slog.Error("mysqluseradapter/GetActiveUserRolesByUserID: error executing Read", "error", err)
		return nil, utils.ErrInternalServer
	}

	if len(entities) == 0 {
		return nil, utils.ErrInternalServer
	}

	if len(entities) > 1 {
		slog.Error("mysqluseradapter.GetActiveUserRolesByUserID: Multiple users found with the same userID", "userID", userID)
		return nil, utils.ErrInternalServer
	}

	role, err = userconverters.UserRoleEntityToDomain(entities[0])
	if err != nil {
		return
	}

	return
}
