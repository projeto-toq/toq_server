package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// PushOptIn handler marks user as opted-in to push notifications (sets opt_status=1).
// Device token is not processed here; it's captured on SignIn.
func (uh *UserHandler) PushOptIn(ctx context.Context, in *pb.PushOptInRequest) (out *pb.PushOptInResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	err = uh.service.PushOptIn(ctx, infos.ID)
	if err != nil {
		return nil, err
	}

	return
}
