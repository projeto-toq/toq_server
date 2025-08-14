package userservices

import (
	"context"

	gcsadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/google_cloud_storage"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) GenerateGCSUploadURL(ctx context.Context, objectName, contentType string) (signedURL string, err error) {
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

	// Validar se é um tipo de foto válido
	validPhotoTypes := map[string]bool{
		gcsadapter.PhotoTypeOriginal: true,
		gcsadapter.PhotoTypeSmall:    true,
		gcsadapter.PhotoTypeMedium:   true,
		gcsadapter.PhotoTypeLarge:    true,
	}

	if !validPhotoTypes[objectName] {
		return "", status.Error(codes.InvalidArgument, "invalid photo type")
	}

	signedURL, err = us.googleCloudService.GeneratePhotoSignedURL(gcsadapter.UsersBucketName, userID, objectName, contentType)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}
