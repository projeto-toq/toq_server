package s3adapter

import (
	"context"

	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateUserFolder prepares user storage namespace in S3.
//
// S3 uses a flat structure with key prefixes instead of directories.
// Prefixes are created automatically when objects are uploaded.
// This function exists for consistency with the CloudStoragePortInterface
// and provides logging for observability purposes.
//
// Parameters:
//   - ctx: Context for logging and tracing
//   - UserID: User's unique identifier
//
// Returns:
//   - error: Always nil (prefixes are created automatically on first object upload)
//
// Note: Legacy .placeholder files are no longer created as they are unnecessary
// for S3 operations and increase storage costs without functional benefit.
func (s *S3Adapter) CreateUserFolder(ctx context.Context, UserID int64) error {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	// S3 automatically creates key prefixes when objects are uploaded
	// No need to create placeholder objects
	logger.Info("adapter.s3.user_folder_ready", "user_id", UserID, "bucket", s.userBucketName, "prefix", UserID)

	return nil
}
