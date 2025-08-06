package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetUserRoles(ctx context.Context, empty *empty.Empty) (response *pb.GetUserRolesResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	roles, err := uh.service.GetUserRolesByUser(ctx, infos.ID)
	if err != nil {
		return
	}

	pbroles := []*pb.UserRole{}

	for _, role := range roles {
		pbrole := pb.UserRole{
			Id:           role.GetID(),
			UserId:       role.GetUserID(),
			BaseRoleId:   role.GetBaseRoleID(),
			Role:         role.GetRole().String(),
			Active:       role.IsActive(),
			Status:       role.GetStatus().String(),
			StatusReason: role.GetStatusReason(),
		}
		pbroles = append(pbroles, &pbrole)
	}

	return &pb.GetUserRolesResponse{
		Roles: pbroles,
	}, nil

}
