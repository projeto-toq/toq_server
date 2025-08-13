package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *UserHandler) GenerateGCSUploadURL(ctx context.Context, req *pb.GenerateGCSUploadURLRequest) (*pb.GenerateGCSUploadURLResponse, error) {
	if req.GetObjectName() == "" || req.GetContentType() == "" {
		return nil, status.Error(codes.InvalidArgument, "object_name and content_type are required")
	}

	signedURL, err := h.service.GenerateGCSUploadURL(ctx, req.ObjectName, req.ContentType)
	if err != nil {
		return nil, err
	}

	return &pb.GenerateGCSUploadURLResponse{
		SignedUrl: signedURL,
	}, nil
}
