package userservices

import (
	"context"
	"fmt"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) GeneratePhotoDownloadURL(ctx context.Context, userID int64) (signedURL string, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	bucketName := fmt.Sprintf("user-%d-bucket", userID)
	objectName := "photo.jpg"

	signedURL, err = us.googleCloudService.GenerateV4GetObjectSignedURL(bucketName, objectName)
	if err != nil {
		return
	}

	return
}
