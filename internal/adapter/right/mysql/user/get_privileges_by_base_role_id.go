package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) GetPrivilegesByBaseRoleID(ctx context.Context, tx *sql.Tx, baseRoleID int64) (privileges []usermodel.PrivilegeInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	query := `SELECT * FROM role_privileges WHERE role_id=?;`

	entities, err := ua.Read(ctx, tx, query, baseRoleID)
	if err != nil {
		slog.Error("mysqluseradapter/GetPrivilegesByBaseRoleID: error executing Read", "error", err)
		return nil, utils.ErrInternalServer
	}

	if len(entities) == 0 {
		return nil, utils.ErrInternalServer
	}

	for _, entity := range entities {
		privilege, err1 := userconverters.PrivilegeEntityToDomain(entity)
		if err1 != nil {
			return nil, err1
		}

		privileges = append(privileges, privilege)
	}

	return
}
