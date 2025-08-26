package userservices

import (
	"context"

	storagemodel "github.com/giulio-alfieri/toq_server/internal/core/model/storage_model"
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
		photoType = string(storagemodel.PhotoOriginal)
	}

	signedURL, err = us.cloudStorageService.GeneratePhotoDownloadURL(userID, storagemodel.PhotoType(photoType))
	if err != nil {
		return
	}

	return
}
