package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) SignIn(ctx context.Context, in *pb.SignInRequest) (response *pb.SignInResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tokens, err := uh.service.SignIn(ctx, in.NationalID, in.Password)
	if err != nil {
		return
	}
	return &pb.SignInResponse{
		Tokens: &pb.Tokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}, nil
}
