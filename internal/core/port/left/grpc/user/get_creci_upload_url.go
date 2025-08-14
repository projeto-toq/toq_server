package grpcuserport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *UserHandler) GetCreciUploadURL(ctx context.Context, req *pb.GetCreciUploadURLRequest) (*pb.GetCreciUploadURLResponse, error) {
	if req.GetDocumentType() == "" || req.GetContentType() == "" {
		return nil, status.Error(codes.InvalidArgument, "document_type and content_type are required")
	}

	signedURL, err := h.service.GetCreciUploadURL(ctx, req.DocumentType, req.ContentType)
	if err != nil {
		return nil, err
	}

	return &pb.GetCreciUploadURLResponse{
		SignedUrl: signedURL,
	}, nil
}
