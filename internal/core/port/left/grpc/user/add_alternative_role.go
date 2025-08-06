package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) AddAlternativeUserRole(ctx context.Context, in *pb.AddAlternativeUserRoleRequest) (response *pb.AddAlternativeUserRoleResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	var alternativeRole usermodel.UserRole

	if infos.Role == usermodel.RoleOwner {
		alternativeRole = usermodel.RoleRealtor
	} else {
		alternativeRole = usermodel.RoleOwner
	}

	err = uh.service.AddAlternativeRole(ctx, infos.ID, alternativeRole, in.GetCreciNumber(), in.GetCreciState(), in.GetCreciValidity())
	if err != nil {
		return
	}

	return &pb.AddAlternativeUserRoleResponse{}, nil

}
