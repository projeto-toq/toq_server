package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	handlervalidators "github.com/giulio-alfieri/toq_server/internal/core/port/left/grpc/validators"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) CreateRealtor(ctx context.Context, in *pb.CreateRealtorRequest) (response *pb.CreateRealtorResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	realtor, err := handlervalidators.CleanAndValidateProfile(ctx, in.GetRealtor())
	if err != nil {
		return
	}

	tokens, err := uh.service.CreateRealtor(ctx, realtor)
	if err != nil {
		return
	}

	return &pb.CreateRealtorResponse{
		Tokens: &pb.Tokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}, nil
}
