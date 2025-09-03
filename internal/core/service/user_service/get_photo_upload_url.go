package userservices

import (
	"context"
	"net/http"

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
		return "", utils.ErrInternalServer
	}

	// Validar se é um tipo de foto válido usando constantes do domínio
	validPhotoTypes := storagemodel.ValidPhotoTypes()
	if !validPhotoTypes[objectName] {
		// 422 para tipo de foto inválido com detalhes do campo
		return "", utils.NewHTTPError(http.StatusUnprocessableEntity, "Invalid photo type", map[string]string{
			"field":   "objectName",
			"message": "Unsupported photo type",
		})
	}

	// Validar content-type permitido (apenas imagens JPEG/PNG)
	allowedContentTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	if !allowedContentTypes[contentType] {
		return "", utils.NewHTTPError(http.StatusUnprocessableEntity, "Invalid content type", map[string]string{
			"field":   "contentType",
			"message": "Only image/jpeg or image/png are allowed",
		})
	}

	signedURL, err = us.cloudStorageService.GeneratePhotoUploadURL(userID, storagemodel.PhotoType(objectName), contentType)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}
