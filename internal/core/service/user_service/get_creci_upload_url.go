package userservices

import (
	"context"

	storagemodel "github.com/giulio-alfieri/toq_server/internal/core/model/storage_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (us *userService) GetCreciUploadURL(ctx context.Context, documentType, contentType string) (signedURL string, err error) {
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

	// Validar se é um tipo de documento CRECI válido usando constantes do domínio
	validCreciTypes := storagemodel.ValidDocumentTypes()

	if !validCreciTypes[documentType] {
		return "", status.Error(codes.InvalidArgument, "invalid CRECI document type")
	}

	// Usar o novo método da interface de storage
	signedURL, err = us.cloudStorageService.GenerateDocumentUploadURL(userID, storagemodel.DocumentType(documentType), contentType)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}
