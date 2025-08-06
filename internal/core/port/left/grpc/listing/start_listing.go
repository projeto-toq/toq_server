package grpclistingport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (lr *ListingHandler) StartListing(ctx context.Context, in *pb.StartListingRequest) (out *pb.StartListingResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	listing, err := lr.service.StartListing(ctx, in.ZipCode, in.Number, globalmodel.PropertyType(in.PropertyType))
	if err != nil {
		return
	}

	out = &pb.StartListingResponse{
		Id: listing.ID(),
	}

	return
}
