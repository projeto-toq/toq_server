package sqsmediaprocessingadapter

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
	"go.opentelemetry.io/otel/trace"
)

// EnqueueJob publishes a payload to the primary queue.
func (a *MediaProcessingQueueAdapter) EnqueueJob(ctx context.Context, payload mediaprocessingmodel.MediaProcessingJobMessage) (string, error) {
	return a.publish(ctx, a.queueURL, payload, "MediaProcessingQueue.EnqueueJob")
}

// EnqueueRetry publishes a payload to the retry queue or falls back to the primary queue.
func (a *MediaProcessingQueueAdapter) EnqueueRetry(ctx context.Context, payload mediaprocessingmodel.MediaProcessingJobMessage) (string, error) {
	target := a.retryQueueURL
	if target == "" {
		target = a.queueURL
	}
	return a.publish(ctx, target, payload, "MediaProcessingQueue.EnqueueRetry")
}

func (a *MediaProcessingQueueAdapter) publish(ctx context.Context, queueURL string, payload mediaprocessingmodel.MediaProcessingJobMessage, operation string) (string, error) {
	ctx = utils.ContextWithLogger(ctx)
	ctx, spanEnd, err := utils.GenerateBusinessTracer(ctx, operation)
	if err != nil {
		return "", derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	if a == nil || a.client == nil || queueURL == "" {
		err := derrors.Infra("media processing queue adapter not configured", nil)
		utils.SetSpanError(ctx, err)
		return "", err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return "", derrors.Infra("failed to marshal queue payload", err)
	}

	traceparent := buildTraceparent(ctx)
	attributes := map[string]types.MessageAttributeValue{
		"ListingIdentityId": stringAttribute(strconv.FormatUint(payload.ListingIdentityID, 10)),
		"JobId":             stringAttribute(strconv.FormatUint(payload.JobID, 10)),
		"RetryCount":        stringAttribute(strconv.FormatUint(uint64(payload.Retry), 10)),
	}
	if traceparent != "" {
		attributes["Traceparent"] = stringAttribute(traceparent)
	}

	input := &sqs.SendMessageInput{
		QueueUrl:          aws.String(queueURL),
		MessageBody:       aws.String(string(body)),
		MessageAttributes: attributes,
	}

	output, err := a.client.SendMessage(ctx, input)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("adapter.sqs.media.send_failed", "queue", queueURL, "listing_identity_id", payload.ListingIdentityID, "error", err)
		return "", derrors.Infra("failed to send SQS message", err)
	}

	logger := utils.LoggerFromContext(ctx)
	logger.Info("adapter.sqs.media.send_success", "queue", queueURL, "listing_identity_id", payload.ListingIdentityID, "message_id", aws.ToString(output.MessageId))

	return aws.ToString(output.MessageId), nil
}

func stringAttribute(value string) types.MessageAttributeValue {
	return types.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String(value),
	}
}

func buildTraceparent(ctx context.Context) string {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if !spanCtx.IsValid() {
		return ""
	}
	return fmt.Sprintf("00-%s-%s-%s", spanCtx.TraceID().String(), spanCtx.SpanID().String(), spanCtx.TraceFlags().String())
}
