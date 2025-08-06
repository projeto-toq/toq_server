package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) DeleteRealtorByID(ctx context.Context, in *pb.DeleteRealtorByIDRequest) (out *pb.DeleteRealtorByIDResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	err = uh.service.DeleteRealtorOfAgency(ctx, infos.ID, in.GetRealtorID())
	if err != nil {
		return
	}

	return
}
