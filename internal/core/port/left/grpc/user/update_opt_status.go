package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// UpdateOptStatus updates the user's push notification opt-in status
func (uh *UserHandler) UpdateOptStatus(ctx context.Context, in *pb.UpdateOptStatusRequest) (out *pb.UpdateOptStatusResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	if err = uh.service.UpdateOptStatus(ctx, infos.ID, in.GetOptIn()); err != nil {
		return nil, err
	}

	out = &pb.UpdateOptStatusResponse{}
	return
}
