package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetUsersByStatus(ctx context.Context, tx *sql.Tx, userRoleStatus usermodel.UserRoleStatus, userRole usermodel.UserRole) (users []usermodel.UserInterface, err error) {

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT user_id FROM user_roles WHERE status = ? AND role = ?;`

	entities, err := ua.Read(ctx, tx, query, userRoleStatus, userRole)
	if err != nil {
		slog.Error("mysqluseradapter/GetUsersByStatus: error executing Read", "error", err)
		return nil, utils.ErrInternalServer
	}

	if len(entities) == 0 {
		return nil, utils.ErrInternalServer
	}

	for _, entity := range entities {
		userID, ok := entity[0].(int64)
		if !ok {
			slog.Error("mysqluseradapter/GetUsersByStatus: error converting user_id to int64", "value", entity[0])
			return nil, utils.ErrInternalServer
		}
		user, err1 := ua.GetUserByID(ctx, tx, userID)
		if err1 != nil {
			return nil, err1
		}
		users = append(users, user)

	}

	return

}
