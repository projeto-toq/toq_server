package userservices

import (
	"context"
	"fmt"

	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	storagemodel "github.com/giulio-alfieri/toq_server/internal/core/model/storage_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

const creciDownloadURLExpirationMinutes = 60

// GetCreciDownloadURLs gera URLs assinadas para download dos documentos CRECI do usuário alvo.
func (us *userService) GetCreciDownloadURLs(ctx context.Context, targetUserID int64) (urls CreciDocumentDownloadURLs, err error) {
	urls = CreciDocumentDownloadURLs{ExpiresInMinutes: creciDownloadURLExpirationMinutes}

	ctx, spanEnd, terr := utils.GenerateTracer(ctx)
	if terr != nil {
		err = utils.InternalError("Failed to generate tracer")
		return
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if targetUserID <= 0 {
		err = utils.ValidationError("id", "User ID must be greater than zero")
		return
	}

	if us.cloudStorageService == nil {
		err = utils.InternalError("Storage service not configured")
		return
	}

	user, uerr := us.GetUserByID(ctx, targetUserID)
	if uerr != nil {
		err = uerr
		return
	}

	activeRole := user.GetActiveRole()
	if activeRole == nil || activeRole.GetRole() == nil {
		err = utils.InternalError("User active role missing")
		return
	}

	if activeRole.GetRole().GetSlug() != permissionmodel.RoleSlugRealtor.String() {
		err = utils.ConflictError("User is not a realtor")
		return
	}

	bucket := us.cloudStorageService.GetBucketConfig().Name
	documents := []storagemodel.DocumentType{storagemodel.DocSelfie, storagemodel.DocFront, storagemodel.DocBack}
	missing := make([]string, 0, len(documents))

	for _, docType := range documents {
		objectPath := fmt.Sprintf("%d/%s", targetUserID, string(docType))
		exists, checkErr := us.cloudStorageService.ObjectExists(ctx, bucket, objectPath)
		if checkErr != nil {
			utils.SetSpanError(ctx, checkErr)
			logger.Error("admin.creci_download.object_exists_error", "user_id", targetUserID, "doc", string(docType), "error", checkErr)
			err = utils.InternalError("Failed to check document existence")
			return
		}
		if !exists {
			missing = append(missing, string(docType))
		}
	}

	if len(missing) > 0 {
		err = utils.NewHTTPErrorWithSource(422, "Missing required documents", map[string]any{"missing": missing})
		return
	}

	selfieURL, sErr := us.cloudStorageService.GenerateDocumentDownloadURL(targetUserID, storagemodel.DocSelfie)
	if sErr != nil {
		utils.SetSpanError(ctx, sErr)
		logger.Error("admin.creci_download.generate_url_error", "user_id", targetUserID, "doc", string(storagemodel.DocSelfie), "error", sErr)
		err = utils.InternalError("Failed to generate download URL")
		return
	}

	frontURL, fErr := us.cloudStorageService.GenerateDocumentDownloadURL(targetUserID, storagemodel.DocFront)
	if fErr != nil {
		utils.SetSpanError(ctx, fErr)
		logger.Error("admin.creci_download.generate_url_error", "user_id", targetUserID, "doc", string(storagemodel.DocFront), "error", fErr)
		err = utils.InternalError("Failed to generate download URL")
		return
	}

	backURL, bErr := us.cloudStorageService.GenerateDocumentDownloadURL(targetUserID, storagemodel.DocBack)
	if bErr != nil {
		utils.SetSpanError(ctx, bErr)
		logger.Error("admin.creci_download.generate_url_error", "user_id", targetUserID, "doc", string(storagemodel.DocBack), "error", bErr)
		err = utils.InternalError("Failed to generate download URL")
		return
	}

	urls.Selfie = selfieURL
	urls.Front = frontURL
	urls.Back = backURL
	return
}
