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

func (ua *UserAdapter) GetBaseRoles(ctx context.Context, tx *sql.Tx) (roles []usermodel.BaseRoleInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM base_roles;`

	entities, err := ua.Read(ctx, tx, query)
	if err != nil {
		slog.Error("mysqluseradapter/GetBaseRoles: error executing Read", "error", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	if len(entities) == 0 {
		return nil, status.Error(codes.NotFound, "Base roles not found")
	}

	for _, entity := range entities {
		role, err1 := userconverters.BaseRoleEntityToDomain(entity)
		if err1 != nil {
			return nil, err1
		}

		privileges, err1 := ua.GetPrivilegesByBaseRoleID(ctx, tx, role.GetID())
		if err1 != nil {
			return nil, err1
		}

		role.SetPrivileges(privileges)
		roles = append(roles, role)
	}

	return
}
