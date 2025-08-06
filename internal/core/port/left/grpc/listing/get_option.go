package grpclistingport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (lr *ListingHandler) GetOptions(ctx context.Context, in *pb.GetOptionsRequest) (out *pb.GetOptionsResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	propertyTypes, err := lr.service.GetOptions(ctx, in.GetZipCode(), in.GetNumber())
	if err != nil {
		return
	}

	out = &pb.GetOptionsResponse{
		PropertyTypes: propertyTypes,
	}

	return
}
