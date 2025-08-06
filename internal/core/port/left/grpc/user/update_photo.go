package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) UpdatePhoto(ctx context.Context, in *pb.UpdatePhotoRequest) (response *pb.UpdatePhotoResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	err = uh.service.UpdatePhoto(ctx, infos.ID, in.GetPhoto())

	return &pb.UpdatePhotoResponse{}, nil
}
