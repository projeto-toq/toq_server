package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// UpdateOptStatus updates the user's push notification opt-in status
func (uh *UserHandler) UpdateOptStatus(ctx context.Context, in *pb.UpdateOptStatusRequest) (out *pb.UpdateOptStatusResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	if err = uh.service.UpdateOptStatus(ctx, in.GetOptIn()); err != nil {
		return nil, err
	}

	out = &pb.UpdateOptStatusResponse{}
	return
}
