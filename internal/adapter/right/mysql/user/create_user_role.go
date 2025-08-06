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

func (ua *UserAdapter) CreateUserRole(ctx context.Context, tx *sql.Tx, role usermodel.UserRoleInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `INSERT INTO user_roles (user_id, base_role_id, role, active, status, status_reason) 
			VALUES (?, ?, ?, ?, ?, ?);`

	entity := userconverters.UserRoleDomainToEntity(role)

	id, err := ua.Create(ctx, tx, sql,
		entity.UserID,
		entity.BaseRoleID,
		entity.Role,
		entity.Active,
		entity.Status,
		entity.StatusReason)
	if err != nil {
		slog.Error("mysqluseradapter/CreateUserRole: error executing Create", "error", err)
		return status.Error(codes.Internal, "Internal server error")

	}

	role.SetID(id)

	return
}
