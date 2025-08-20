package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) ResendPhoneChangeCode(ctx context.Context, in *pb.ResendPhoneChangeCodeRequest) (response *pb.ResendPhoneChangeCodeResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	code, err := uh.service.ResendPhoneChangeCode(ctx, infos.ID)
	if err != nil {
		return
	}

	return &pb.ResendPhoneChangeCodeResponse{
		Code: code,
	}, nil
}
