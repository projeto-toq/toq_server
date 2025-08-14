package userservices

import (
	"context"

	gcsadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/google_cloud_storage"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) GetCreciUploadURL(ctx context.Context, documentType, contentType string) (signedURL string, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil {
		return
	}

	if us.googleCloudService == nil {
		return "", status.Error(codes.Unimplemented, "GCS service is not configured")
	}

	// Validar se é um tipo de documento CRECI válido
	validCreciTypes := map[string]bool{
		gcsadapter.CreciDocumentSelfie: true,
		gcsadapter.CreciDocumentFront:  true,
		gcsadapter.CreciDocumentBack:   true,
	}

	if !validCreciTypes[documentType] {
		return "", status.Error(codes.InvalidArgument, "invalid CRECI document type")
	}

	// Usar o mesmo método de geração de URL, mas com documentType
	signedURL, err = us.googleCloudService.GeneratePhotoSignedURL(gcsadapter.UsersBucketName, userID, documentType, contentType)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}
