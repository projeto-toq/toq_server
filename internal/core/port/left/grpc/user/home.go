package grpcuserport

import (
	"context"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GoHome(ctx context.Context, empty *empty.Empty) (response *pb.GoHomeResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	user, err := uh.service.Home(ctx, infos.ID)
	if err != nil {
		return nil, err
	}

	return &pb.GoHomeResponse{
		Message: fmt.Sprintf("Welcome %s. Seu Role é %s, e seu perfil está %v", user.GetNickName(), infos.Role, infos.ProfileStatus),
	}, nil
}
