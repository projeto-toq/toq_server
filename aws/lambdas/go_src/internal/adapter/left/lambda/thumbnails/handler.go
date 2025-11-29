package thumbnails

import (
	"context"
	"log/slog"

	"github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/service/image_processing"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// Handler manages the Lambda execution flow
type Handler struct {
	service *imageprocessing.ThumbnailService
	logger  *slog.Logger
}

// NewHandler creates a new Lambda handler
func NewHandler(service *imageprocessing.ThumbnailService, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type ThumbnailOutput struct {
	Status          string                          `json:"status"`
	GeneratedAssets []mediaprocessingmodel.JobAsset `json:"generatedAssets"`
}

// HandleRequest processes the Step Function payload
//
// @Summary     Process image thumbnails
// @Description Receives a batch of assets, filters for photos, and generates thumbnails (Large, Medium, Small, Tiny).
//
//	Returns the list of generated assets.
func (h *Handler) HandleRequest(ctx context.Context, event mediaprocessingmodel.StepFunctionPayload) (mediaprocessingmodel.LambdaResponse, error) {
	h.logger.Info("Thumbnails Lambda started", "batch_id", event.BatchID, "assets_to_process", len(event.ValidAssets))

	allGeneratedAssets := make([]mediaprocessingmodel.JobAsset, 0)

	for i, asset := range event.ValidAssets {
		h.logger.Info("Inspecting asset", "index", i, "key", asset.Key, "type", asset.Type)

		generated, err := h.service.ProcessAsset(ctx, event.ListingID, asset)
		if err != nil {
			h.logger.Error("Failed to process asset", "key", asset.Key, "error", err)
			// We continue processing other assets even if one fails
			continue
		}

		if generated != nil {
			allGeneratedAssets = append(allGeneratedAssets, generated...)
		}
	}

	h.logger.Info("Thumbnails Lambda finished", "generated_count", len(allGeneratedAssets))

	return mediaprocessingmodel.LambdaResponse{
		Body: ThumbnailOutput{
			Status:          "SUCCESS",
			GeneratedAssets: allGeneratedAssets,
		},
	}, nil
}
