package grpclistingport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (lr *ListingHandler) GetBaseFeatures(ctx context.Context, in *pb.GetBaseFeaturesRequest) (out *pb.GetBaseFeaturesResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	features, err := lr.service.GetBaseFeatures(ctx)
	if err != nil {
		return
	}

	response := make([]*pb.BaseFeature, 0, len(features))

	for _, feature := range features {
		response = append(response, &pb.BaseFeature{
			Id:          feature.ID(),
			Feature:     feature.Feature(),
			Description: feature.Description(),
		})
	}

	out = &pb.GetBaseFeaturesResponse{
		Features: response,
	}

	return
}
