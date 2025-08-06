package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) InviteRealtor(ctx context.Context, in *pb.InviteRealtorRequest) (out *pb.InviteRealtorResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	err = uh.service.InviteRealtor(ctx, in.GetPhoneNumber())
	if err != nil {
		return
	}

	return &pb.InviteRealtorResponse{}, nil

}
