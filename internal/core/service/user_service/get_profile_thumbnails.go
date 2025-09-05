package userservices

import (
	"context"
	"log/slog"

	storagemodel "github.com/giulio-alfieri/toq_server/internal/core/model/storage_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) GetProfileThumbnails(ctx context.Context) (thumbnails usermodel.ProfileThumbnails, err error) {
	// Obter o ID do usuário do contexto (SSOT)
	userID, err := us.globalService.GetUserIDFromContext(ctx)
	if err != nil || userID == 0 {
		return thumbnails, utils.InternalError("Failed to get user from context")
	}
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	// Gerar URLs assinadas para cada tipo de foto usando a nova interface
	originalURL, err := us.cloudStorageService.GeneratePhotoDownloadURL(userID, storagemodel.PhotoOriginal)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.get_profile_thumbnails.generate_original_url_error", "error", err, "user_id", userID)
		return
	}

	smallURL, err := us.cloudStorageService.GeneratePhotoDownloadURL(userID, storagemodel.PhotoSmall)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.get_profile_thumbnails.generate_small_url_error", "error", err, "user_id", userID)
		return
	}

	mediumURL, err := us.cloudStorageService.GeneratePhotoDownloadURL(userID, storagemodel.PhotoMedium)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.get_profile_thumbnails.generate_medium_url_error", "error", err, "user_id", userID)
		return
	}

	largeURL, err := us.cloudStorageService.GeneratePhotoDownloadURL(userID, storagemodel.PhotoLarge)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.get_profile_thumbnails.generate_large_url_error", "error", err, "user_id", userID)
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
