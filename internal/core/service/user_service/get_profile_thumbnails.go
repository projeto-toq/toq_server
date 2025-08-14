package userservices

import (
	"context"

	gcsadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/google_cloud_storage"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) GetProfileThumbnails(ctx context.Context, userID int64) (thumbnails usermodel.ProfileThumbnails, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Gerar URLs assinadas para cada tipo de foto
	originalURL, err := us.googleCloudService.GeneratePhotoDownloadURL(gcsadapter.UsersBucketName, userID, gcsadapter.PhotoTypeOriginal)
	if err != nil {
		return
	}

	smallURL, err := us.googleCloudService.GeneratePhotoDownloadURL(gcsadapter.UsersBucketName, userID, gcsadapter.PhotoTypeSmall)
	if err != nil {
		return
	}

	mediumURL, err := us.googleCloudService.GeneratePhotoDownloadURL(gcsadapter.UsersBucketName, userID, gcsadapter.PhotoTypeMedium)
	if err != nil {
		return
	}

	largeURL, err := us.googleCloudService.GeneratePhotoDownloadURL(gcsadapter.UsersBucketName, userID, gcsadapter.PhotoTypeLarge)
	if err != nil {
		return
	}

	thumbnails = usermodel.ProfileThumbnails{
		OriginalURL: originalURL,
		SmallURL:    smallURL,
		MediumURL:   mediumURL,
		LargeURL:    largeURL,
	}

	return
}
