package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
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
		return nil, utils.ErrInternalServer
	}

	if len(entities) == 0 {
		return nil, utils.ErrInternalServer
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
