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

func (ua *UserAdapter) GetBaseRoleByRole(ctx context.Context, tx *sql.Tx, roleName usermodel.UserRole) (role usermodel.BaseRoleInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM base_roles WHERE role=?;`

	entities, err := ua.Read(ctx, tx, query, roleName)
	if err != nil {
		slog.Error("mysqluseradapter/GetBaseRoleByRole: error executing Read", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if len(entities) == 0 {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	if len(entities) > 1 {
		slog.Error("UserAdapter.GetBaseRoleByRole: Multiple users found with the same role", "roleName", roleName)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	role, err = userconverters.BaseRoleEntityToDomain(entities[0])
	if err != nil {
		return
	}

	privileges, err := ua.GetPrivilegesByBaseRoleID(ctx, tx, role.GetID())
	if err != nil {
		return
	}

	role.SetPrivileges(privileges)

	return
}
