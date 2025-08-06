package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetOnboardingStatus(ctx context.Context, in *pb.GetOnboardingStatusRequest) (response *pb.GetOnboardingStatusResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	status, reason, err := uh.service.GetOnboardingStatus(ctx, infos.ID)
	if err != nil {
		return nil, err
	}
	response = &pb.GetOnboardingStatusResponse{
		Status: status,
		Reason: reason,
	}

	return
}
