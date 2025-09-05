package userservices

import (
	"context"
	"log/slog"

	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (us *userService) BatchUpdateLastActivity(ctx context.Context, userIDs []int64, timestamps []int64) (err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	if len(userIDs) != len(timestamps) {
		return utils.ValidationError("timestamps", "userIDs and timestamps length mismatch")
	}

	// Call repository batch update method
	err = us.repo.BatchUpdateUserLastActivity(ctx, userIDs, timestamps)
	if err != nil {
		utils.SetSpanError(ctx, err)
		slog.Error("user.batch_update_last_activity.repo_error", "error", err, "count", len(userIDs))
		return utils.InternalError("Failed to batch update last activity")
	}

	return
}
