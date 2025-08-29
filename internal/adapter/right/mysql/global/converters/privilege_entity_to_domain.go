package globalconverters

import (
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func PrivilegeEntityToDomain(entity []any) (privilege usermodel.PrivilegeInterface, err error) {
	privilege = usermodel.NewPrivilege()

	id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "value", entity[0])
		return nil, utils.ErrInternalServer
	}
	privilege.SetID(id)

	roleID, ok := entity[1].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "value", entity[1])
		return nil, utils.ErrInternalServer
	}
	privilege.SetRoleID(roleID)

	service, ok := entity[2].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "value", entity[2])
		return nil, utils.ErrInternalServer
	}
	privilege.SetService(usermodel.GRPCService(service))

	method, ok := entity[3].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "value", entity[3])
		return nil, utils.ErrInternalServer
	}
	privilege.SetMethod(uint8(method))

	allowed, ok := entity[4].(int64)
	if !ok {
		slog.Error("Error converting active to int64", "value", entity[4])
		return nil, utils.ErrInternalServer
	}
	privilege.SetAllowed(allowed == 1)

	return
}
