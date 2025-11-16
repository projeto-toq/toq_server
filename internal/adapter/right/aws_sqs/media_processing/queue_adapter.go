package sqsmediaprocessingadapter

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	mediaprocessingqueue "github.com/projeto-toq/toq_server/internal/core/port/right/queue/mediaprocessingqueue"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// MediaProcessingQueueAdapter publishes media processing jobs to AWS SQS.
type MediaProcessingQueueAdapter struct {
	client        *sqs.Client
	queueURL      string
	retryQueueURL string
	region        string
}

// NewMediaProcessingQueueAdapter configures the adapter using environment data.
func NewMediaProcessingQueueAdapter(ctx context.Context, env *globalmodel.Environment) (*MediaProcessingQueueAdapter, error) {
	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	queueURL := env.MediaProcessing.Queue.URL
	if queueURL == "" {
		logger.Warn("adapter.sqs.media.queue_url_missing")
		return nil, nil
	}

	region := env.MediaProcessing.Queue.Region
	if region == "" {
		region = env.S3.Region
	}
	if region == "" {
		return nil, derrors.Validation("media processing queue region not configured", nil)
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		logger.Error("adapter.sqs.media.config_error", "region", region, "error", err)
		return nil, derrors.Infra("failed to load AWS config", err)
	}

	client := sqs.NewFromConfig(cfg)
	logger.Info("adapter.sqs.media.created", "region", region, "queue", queueURL)

	return &MediaProcessingQueueAdapter{
		client:        client,
		queueURL:      queueURL,
		retryQueueURL: env.MediaProcessing.Queue.RetryURL,
		region:        region,
	}, nil
}

var _ mediaprocessingqueue.QueuePortInterface = (*MediaProcessingQueueAdapter)(nil)
