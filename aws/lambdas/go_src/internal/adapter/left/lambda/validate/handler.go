package validate

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/service/validate"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

type Handler struct {
	service *validate.ValidateService
	logger  *slog.Logger
}

func NewHandler(service *validate.ValidateService, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) HandleRequest(ctx context.Context, rawEvent json.RawMessage) (mediaprocessingmodel.StepFunctionPayload, error) {
	h.logger.Info("Validate Lambda received raw event", "raw_event", string(rawEvent))

	// Try to parse as SQS Event
	var sqsEvent events.SQSEvent
	if err := json.Unmarshal(rawEvent, &sqsEvent); err == nil && len(sqsEvent.Records) > 0 && sqsEvent.Records[0].EventSource == "aws:sqs" {
		h.logger.Info("Detected SQS Event", "record_count", len(sqsEvent.Records))
		if err := h.service.ProcessSQSEvent(ctx, sqsEvent.Records); err != nil {
			return mediaprocessingmodel.StepFunctionPayload{}, err
		}
		return mediaprocessingmodel.StepFunctionPayload{}, nil
	}

	var event mediaprocessingmodel.StepFunctionPayload
	if err := json.Unmarshal(rawEvent, &event); err != nil {
		h.logger.Error("Failed to unmarshal event", "error", err)
		return mediaprocessingmodel.StepFunctionPayload{}, err
	}

	return h.service.ValidateAssets(ctx, event)
}
