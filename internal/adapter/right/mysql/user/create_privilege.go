package mysqluseradapter

import (
	"context"
	"database/sql"
	"log/slog"

	userconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/user/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ua *UserAdapter) CreatePrivileges(ctx context.Context, tx *sql.Tx, privileges []usermodel.PrivilegeInterface, roleID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	sql := `INSERT INTO role_privileges (role_id, service, method, allowed) VALUES (?, ?, ?, ?);`

	for _, privilege := range privileges {
		entity := userconverters.PrivilegeDomainToEntity(privilege, roleID)

		id, err1 := ua.Create(ctx, tx, sql, entity.RoleID, entity.Service, entity.Method, entity.Allowed)
		if err1 != nil {
			slog.Error("mysqluseradapter/CreatePrivileges: error executing Create", "error", err1)
			return utils.ErrInternalServer
		}

		privilege.SetID(id)
	}

	return
}
