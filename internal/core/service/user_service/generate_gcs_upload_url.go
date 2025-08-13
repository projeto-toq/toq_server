package userservices

import (
	"context"
	"fmt"

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

	bucketName := fmt.Sprintf("user-%d-bucket", userID)

	signedURL, err = us.googleCloudService.GenerateV4PutObjectSignedURL(bucketName, objectName, contentType)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}
