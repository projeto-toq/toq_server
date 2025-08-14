package userservices

import (
	"context"

	gcsadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/google_cloud_storage"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) GeneratePhotoDownloadURL(ctx context.Context, userID int64, photoType string) (signedURL string, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Se photoType não for especificado, usar o original
	if photoType == "" {
		photoType = gcsadapter.PhotoTypeOriginal
	}

	signedURL, err = us.googleCloudService.GeneratePhotoDownloadURL(gcsadapter.UsersBucketName, userID, photoType)
	if err != nil {
		return
	}

	return
}
