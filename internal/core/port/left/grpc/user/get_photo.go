package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) GetPhoto(ctx context.Context, in *pb.GetPhotoRequest) (response *pb.GetPhotoResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	photo, err := uh.service.GetPhoto(ctx, infos.ID)
	if err != nil {
		return
	}

	return &pb.GetPhotoResponse{
		Photo: photo,
	}, nil
}
