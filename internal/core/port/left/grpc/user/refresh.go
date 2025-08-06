package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) RefreshToken(ctx context.Context, in *pb.RefreshTokenRequest) (response *pb.RefreshTokenResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	tokens, err := uh.service.RefreshTokens(ctx, in.RefreshToken)

	if err != nil {
		return
	}

	return &pb.RefreshTokenResponse{
		Tokens: &pb.Tokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}, nil
}
