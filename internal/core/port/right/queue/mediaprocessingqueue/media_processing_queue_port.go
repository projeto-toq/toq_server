package mediaprocessingqueue

import (
	"context"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// QueuePortInterface defines the contract with the async pipeline used for media processing jobs.
type QueuePortInterface interface {
	EnqueueJob(ctx context.Context, payload mediaprocessingmodel.MediaProcessingJobMessage) (string, error)
	EnqueueRetry(ctx context.Context, payload mediaprocessingmodel.MediaProcessingJobMessage) (string, error)
	DecodeMessage(ctx context.Context, rawBody string) (mediaprocessingmodel.MediaProcessingCallback, error)
	Acknowledge(ctx context.Context, receiptHandle string) error
}
