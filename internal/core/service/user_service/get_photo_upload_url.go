package userservices

import (
	"context"

	storagemodel "github.com/giulio-alfieri/toq_server/internal/core/model/storage_model"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) GetPhotoUploadURL(ctx context.Context, objectName, contentType string) (signedURL string, err error) {
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

	// Validar se é um tipo de foto válido usando constantes do domínio
	validPhotoTypes := storagemodel.ValidPhotoTypes()
	if !validPhotoTypes[objectName] {
		// 422 para tipo de foto inválido com detalhes do campo
		return "", utils.ValidationError("objectName", "Unsupported photo type")
	}

	// Validar content-type permitido via util compartilhado
	if !utils.IsAllowedImageContentType(contentType) {
		return "", utils.ValidationError("contentType", "Only image/jpeg or image/png are allowed")
	}

	signedURL, err = us.cloudStorageService.GeneratePhotoUploadURL(userID, storagemodel.PhotoType(objectName), contentType)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return "", err
	}

	return signedURL, nil
}
