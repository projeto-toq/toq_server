package userconverters

import (
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	
	
"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func UserRoleEntityToDomain(entity []any) (role usermodel.UserRoleInterface, err error) {
	role = usermodel.NewUserRole()

	id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "value", entity[0])
		return nil, utils.ErrInternalServer
	}
	role.SetID(id)

	user_id, ok := entity[1].(int64)
	if !ok {
		slog.Error("Error converting user_id to int64", "value", entity[1])
		return nil, utils.ErrInternalServer
	}
	role.SetUserID(user_id)

	base_role_id, ok := entity[2].(int64)
	if !ok {
		slog.Error("Error converting base_role_id to int64", "value", entity[2])
		return nil, utils.ErrInternalServer
	}
	role.SetBaseRoleID(base_role_id)

	entity_role, ok := entity[3].(int64)
	if !ok {
		slog.Error("Error converting national_id to int64", "value", entity[3])
		return nil, utils.ErrInternalServer
	}
	role.SetRole(usermodel.UserRole(entity_role))

	active, ok := entity[4].(int64)
	if !ok {
		slog.Error("Error converting active to int64", "value", entity[4])
		return nil, utils.ErrInternalServer
	}
	role.SetActive(active == 1)

	role_status, ok := entity[5].(int64)
	if !ok {
		slog.Error("Error converting role_status to int64", "value", entity[5])
		return nil, utils.ErrInternalServer
	}
	role.SetStatus(usermodel.UserRoleStatus(role_status))

	if entity[6] != nil {
		status_reason, ok := entity[6].([]byte)
		if !ok {
			slog.Error("Error converting status_reason to []byte", "value", entity[6])
			return nil, utils.ErrInternalServer
		}
		role.SetStatusReason(string(status_reason))
	}

	return
}
