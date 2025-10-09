package userservices

import (
	"context"

	storagemodel "github.com/projeto-toq/toq_server/internal/core/model/storage_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetPhotoDownloadURL generates a signed URL for a single photo variant
func (us *userService) GetPhotoDownloadURL(ctx context.Context, variant string) (signedURL string, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil {
		return "", err
	}

	if us.cloudStorageService == nil {
		return "", utils.InternalError("Storage service not configured")
	}

	var photoType storagemodel.PhotoType
	switch variant {
	case "original":
		photoType = storagemodel.PhotoOriginal
	case "small":
		photoType = storagemodel.PhotoSmall
	case "medium":
		photoType = storagemodel.PhotoMedium
	case "large":
		photoType = storagemodel.PhotoLarge
	default:
		return "", utils.ValidationError("variant", "Unsupported photo variant")
	}

	signedURL, err = us.cloudStorageService.GeneratePhotoDownloadURL(userID, photoType)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return "", err
	}
	return signedURL, nil
}
