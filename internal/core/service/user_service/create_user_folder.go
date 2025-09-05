package userservices

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) CreateUserFolder(ctx context.Context, userID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("Failed to generate tracer")
	}
	defer spanEnd()

	err = us.cloudStorageService.CreateUserFolder(ctx, userID)
	if err != nil {
		// Provider de storage é infra
		utils.SetSpanError(ctx, err)
		slog.Error("user.create_user_folder.provider_error", "user_id", userID, "err", err)
		return utils.InternalError("Failed to create user folder")
	}

	return nil
}
