package userservices

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) DeleteUserFolder(ctx context.Context, userID int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	err = us.googleCloudService.DeleteUserFolder(ctx, userID)
	if err != nil {
		return
	}

	return
}
