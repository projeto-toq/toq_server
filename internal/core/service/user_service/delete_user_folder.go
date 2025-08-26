package userservices

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) DeleteUserFolder(ctx context.Context, userID int64) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	return us.cloudStorageService.DeleteUserFolder(ctx, userID)
}
