package fcmadapter

import (
	"context"

	"firebase.google.com/go/messaging"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// SendMultipleMessages sends a notification to multiple device tokens
// Firebase supports up to 500 tokens per batch request
func (f *FCMAdapter) SendMultipleMessages(ctx context.Context, message globalmodel.Notification, deviceTokens []string) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	if len(deviceTokens) == 0 {
		logger.Warn("fcm.send_multiple.empty_tokens")
		return nil
	}

	// Firebase allows up to 500 tokens per batch
	const maxBatchSize = 500

	for i := 0; i < len(deviceTokens); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(deviceTokens) {
			end = len(deviceTokens)
		}

		batch := deviceTokens[i:end]
		if err := f.sendBatch(ctx, message, batch); err != nil {
			logger.Error("fcm.send_multiple.batch_error", "batch_start", i, "batch_end", end, "error", err)
			return err
		}

		logger.Info("fcm.send_multiple.batch_success", "tokens_count", len(batch), "batch_start", i)
	}

	return nil
}

func (f *FCMAdapter) sendBatch(ctx context.Context, message globalmodel.Notification, tokens []string) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	multicastMessage := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    message.Title,
			Body:     message.Body,
			ImageURL: message.Icon,
		},
		Tokens: tokens,
	}

	response, err := f.client.SendMulticast(ctx, multicastMessage)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("fcm.send_multiple.multicast_error", "error", err)
		return err
	}

	// Log results
	logger.Info("fcm.send_multiple.multicast_success",
		"success_count", response.SuccessCount,
		"failure_count", response.FailureCount,
		"total_tokens", len(tokens))

	// Log failed tokens for debugging
	if response.FailureCount > 0 {
		for i, resp := range response.Responses {
			if !resp.Success {
				logger.Warn("fcm.send_multiple.token_error",
					"token_index", i,
					"error", resp.Error)
			}
		}
	}

	return nil
}
