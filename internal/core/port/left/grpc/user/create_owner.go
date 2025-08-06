package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	handlervalidators "github.com/giulio-alfieri/toq_server/internal/core/port/left/grpc/validators"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) CreateOwner(ctx context.Context, in *pb.CreateOwnerRequest) (response *pb.CreateOwnerResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	owner, err := handlervalidators.CleanAndValidateProfile(ctx, in.GetOwner())
	if err != nil {
		return
	}

	tokens, err := uh.service.CreateOwner(ctx, owner)
	if err != nil {
		return
	}

	return &pb.CreateOwnerResponse{
		Tokens: &pb.Tokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}, nil
}
