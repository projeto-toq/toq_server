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

func (ua *UserAdapter) CreateBaseRole(ctx context.Context, tx *sql.Tx, role usermodel.BaseRoleInterface) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `INSERT INTO base_roles (role, name) VALUES (?, ?);`

	entity := userconverters.BaseRoleDomainToEntity(role)

	id, err := ua.Create(ctx, tx, sql, entity.Role, entity.Name)
	if err != nil {
		slog.Error("mysqluseradapter/CreateBaseRole: error executing Create", "error", err)
		return status.Error(codes.Internal, "Internal server error")
	}

	role.SetID(id)

	err = ua.CreatePrivileges(ctx, tx, role.GetPrivileges(), role.GetID())
	if err != nil {
		return
	}

	return
}
