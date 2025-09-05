package userservices

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) DeleteUserFolder(ctx context.Context, userID int64) error {
	// No span here: follow guideline to keep spans only in public service methods.
	if us.cloudStorageService == nil {
		return nil
	}
	if err := us.cloudStorageService.DeleteUserFolder(ctx, userID); err != nil {
		// Mark current span error and log infra failure; caller wraps/masks as needed.
		utils.SetSpanError(ctx, err)
		slog.Error("user.delete_user_folder.provider_error", "error", err, "user_id", userID)
		return err
	}
	return nil
}
