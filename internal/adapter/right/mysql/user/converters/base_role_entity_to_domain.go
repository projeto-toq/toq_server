package userconverters

import (
	"log/slog"

	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func BaseRoleEntityToDomain(entity []any) (domain usermodel.BaseRoleInterface, err error) {
	domain = usermodel.NewBaseRole()

	id, ok := entity[0].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "value", entity[0])
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	domain.SetID(id)

	role, ok := entity[1].(int64)
	if !ok {
		slog.Error("Error converting ID to int64", "value", entity[1])
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	domain.SetRole(usermodel.UserRole(role))

	if entity[2] != nil {
		name, ok := entity[2].([]byte)
		if !ok {
			slog.Error("Error converting status_reason to []byte", "value", entity[2])
			return nil, status.Error(codes.Internal, "Internal server error")
		}
		domain.SetName(string(name))
	}

	return
}
