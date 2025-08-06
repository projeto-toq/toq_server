package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/golang/protobuf/ptypes/empty"
)

func (uh *UserHandler) DeleteAgencyOfRealtor(ctx context.Context, empty *empty.Empty) (out *pb.DeleteAgencyOfRealtorResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	err = uh.service.DeleteAgencyOfRealtor(ctx, infos.ID)
	if err != nil {
		return
	}

	return
}
