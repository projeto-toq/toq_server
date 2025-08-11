package fcmadapter

import (
	"context"
	"log/slog"

	"firebase.google.com/go/messaging"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SendMultipleMessages sends a notification to multiple device tokens
// Firebase supports up to 500 tokens per batch request
func (f *FCMAdapter) SendMultipleMessages(ctx context.Context, message globalmodel.Notification, deviceTokens []string) error {
	if len(deviceTokens) == 0 {
		slog.Warn("no device tokens provided for multiple message send")
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
			slog.Error("failed to send batch", "batch_start", i, "batch_end", end, "error", err)
			return err
		}

		slog.Info("batch sent successfully", "tokens_count", len(batch), "batch_start", i)
	}

	return nil
}

func (f *FCMAdapter) sendBatch(ctx context.Context, message globalmodel.Notification, tokens []string) error {
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
		slog.Error("failed to send multicast message", "error", err)
		return status.Error(codes.Internal, "failed to send push notification")
	}

	// Log results
	slog.Info("multicast message sent",
		"success_count", response.SuccessCount,
		"failure_count", response.FailureCount,
		"total_tokens", len(tokens))

	// Log failed tokens for debugging
	if response.FailureCount > 0 {
		for i, resp := range response.Responses {
			if !resp.Success {
				slog.Warn("failed to send to token",
					"token_index", i,
					"error", resp.Error)
			}
		}
	}

	return nil
}
