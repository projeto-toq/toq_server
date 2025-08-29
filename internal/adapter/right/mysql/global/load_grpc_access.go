package mysqlglobaladapter

import (
	"context"
	"database/sql"
	"log/slog"

	globalconverters "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql/global/converters"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (ga *GlobalAdapter) LoadGRPCAccess(ctx context.Context, tx *sql.Tx, service usermodel.GRPCService, method uint8, role usermodel.UserRole) (privilege usermodel.PrivilegeInterface, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	slog.Debug("LoadGRPCAccess called", "service", service, "method", method, "role", role)

	query := `SELECT rp.*
	FROM role_privileges rp
	JOIN base_roles br ON br.id = rp.role_id
	WHERE br.role = ?;`

	entities, err := ga.Read(ctx, tx, query, uint8(role))
	if err != nil {
		slog.Error("mysqluseradapter/GetPrivilegesByBaseRoleID: error executing Read", "error", err)
		return nil, utils.ErrInternalServer
	}

	slog.Debug("LoadGRPCAccess query result", "role", role, "entities_count", len(entities))

	if len(entities) == 0 {
		slog.Debug("No privileges found for role", "role", role)
		return nil, utils.ErrInternalServer
	}

	for _, entity := range entities {
		privilege, err = globalconverters.PrivilegeEntityToDomain(entity)
		if err != nil {
			return
		}
		slog.Debug("Checking privilege", "service", privilege.Service(), "method", privilege.Method(), "allowed", privilege.Allowed(), "target_service", service, "target_method", method)
		if privilege.Service() == service && privilege.Method() == method {
			slog.Debug("Found matching privilege", "allowed", privilege.Allowed())
			return
		}
	}
	slog.Debug("No matching privilege found", "service", service, "method", method, "role", role)
	return
}
