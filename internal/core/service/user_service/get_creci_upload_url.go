package userservices

import (
	"context"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	storagemodel "github.com/projeto-toq/toq_server/internal/core/model/storage_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// GetCreciUploadURL generates a signed URL to upload CRECI documents (selfie/front/back)
// Business rules:
// - Only users with role slug "realtor" are allowed (if claim is present)
// - documentType must be one of storagemodel.ValidDocumentTypes()
// - contentType restricted via utils.IsAllowedImageContentType
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
		return "", utils.InternalError("Storage service not configured")
	}

	// Optional role check (defensive): only realtor can upload CRECI docs
	if ui, ok := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos); ok {
		if ui.RoleSlug != "realtor" { // slug específico
			return "", utils.AuthorizationError("Only realtor role can upload CRECI documents")
		}
	}

	// Validate document type
	if !storagemodel.ValidDocumentTypes()[documentType] {
		return "", utils.ValidationError("documentType", "Unsupported document type")
	}

	// Validate content type
	if !utils.IsAllowedImageContentType(contentType) {
		return "", utils.ValidationError("contentType", "Only image/jpeg or image/png are allowed")
	}

	signedURL, err = us.cloudStorageService.GenerateDocumentUploadURL(userID, storagemodel.DocumentType(documentType), contentType)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return "", utils.InternalError("Failed to generate upload URL")
	}

	return signedURL, nil
}
