package userservices

import (
	"context"

	storagemodel "github.com/projeto-toq/toq_server/internal/core/model/storage_model"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (us *userService) GetPhotoUploadURL(ctx context.Context, variant, contentType string) (signedURL string, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil {
		return
	}

	if us.cloudStorageService == nil {
		return "", utils.InternalError("Storage service not configured")
	}

	// Validar variant aceito e mapear para PhotoType
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

	// Validar content-type permitido via util compartilhado
	if !utils.IsAllowedImageContentType(contentType) {
		return "", utils.ValidationError("contentType", "Only image/jpeg or image/png are allowed")
	}

	signedURL, err = us.cloudStorageService.GeneratePhotoUploadURL(userID, photoType, contentType)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return "", err
	}

	return signedURL, nil
}
