package grpclistingport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (lr *ListingHandler) RejectVisiting(ctx context.Context, in *pb.RejectVistingRequest) (out *pb.RejectVistingResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	_ = infos

	// err = uh.service.AcceptInvitation(ctx, infos.ID)

	return
}
