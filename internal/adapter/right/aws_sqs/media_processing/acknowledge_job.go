package sqsmediaprocessingadapter

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// DecodeMessage converts the Step Functions callback payload carried through SQS into the domain struct.
func (a *MediaProcessingQueueAdapter) DecodeMessage(ctx context.Context, rawBody string) (mediaprocessingmodel.MediaProcessingCallback, error) {
	ctx = utils.ContextWithLogger(ctx)
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "MediaProcessingQueue.DecodeMessage")
	if err != nil {
		return mediaprocessingmodel.MediaProcessingCallback{}, derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if rawBody == "" {
		return mediaprocessingmodel.MediaProcessingCallback{}, derrors.Validation("empty callback payload", nil)
	}

	var callback mediaprocessingmodel.MediaProcessingCallback
	if err := json.Unmarshal([]byte(rawBody), &callback); err != nil {
		utils.SetSpanError(ctx, err)
		return mediaprocessingmodel.MediaProcessingCallback{}, derrors.Validation("invalid callback payload", map[string]string{"error": err.Error()})
	}
	callback.RawBody = rawBody

	return callback, nil
}

// Acknowledge deletes a processed message from the queue.
func (a *MediaProcessingQueueAdapter) Acknowledge(ctx context.Context, receiptHandle string) error {
	ctx = utils.ContextWithLogger(ctx)
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, "MediaProcessingQueue.Acknowledge")
	if err != nil {
		return derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if a == nil || a.client == nil || a.queueURL == "" {
		err := derrors.Infra("media processing queue adapter not configured", nil)
		utils.SetSpanError(ctx, err)
		return err
	}

	if receiptHandle == "" {
		err := derrors.Validation("receipt handle is required", nil)
		utils.SetSpanError(ctx, err)
		return err
	}

	_, err = a.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(a.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("adapter.sqs.media.ack_failed", "queue", a.queueURL, "error", err)
		return derrors.Infra("failed to delete SQS message", err)
	}

	logger := utils.LoggerFromContext(ctx)
	logger.Info("adapter.sqs.media.ack_success", "queue", a.queueURL)
	return nil
}
