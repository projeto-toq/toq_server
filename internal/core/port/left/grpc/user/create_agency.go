package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	handlervalidators "github.com/giulio-alfieri/toq_server/internal/core/port/left/grpc/validators"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (uh *UserHandler) CreateAgency(ctx context.Context, in *pb.CreateAgencyRequest) (response *pb.CreateAgencyResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	agency, err := handlervalidators.CleanAndValidateProfile(ctx, in.GetAgency())
	if err != nil {
		return
	}

	tokens, err := uh.service.CreateAgency(ctx, agency)
	if err != nil {
		return
	}

	return &pb.CreateAgencyResponse{
		Tokens: &pb.Tokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}, nil
}
