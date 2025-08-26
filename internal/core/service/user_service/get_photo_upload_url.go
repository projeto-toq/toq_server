package userservices

import (
	"context"

	storagemodel "github.com/giulio-alfieri/toq_server/internal/core/model/storage_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return "", status.Error(codes.Unimplemented, "Cloud storage service is not configured")
	}

	// Validar se é um tipo de foto válido usando constantes do domínio
	validPhotoTypes := storagemodel.ValidPhotoTypes()

	if !validPhotoTypes[objectName] {
		return "", status.Error(codes.InvalidArgument, "invalid photo type")
	}

	signedURL, err = us.cloudStorageService.GeneratePhotoUploadURL(userID, storagemodel.PhotoType(objectName), contentType)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}
