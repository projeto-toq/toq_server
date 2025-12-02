package thumbnails

import (
	"context"
	"log/slog"
	"os"

	imageprocessing "github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/service/image_processing"
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
	Errors          []ThumbnailError                `json:"errors"`
}

// ThumbnailError captures failures during thumbnail processing so the pipeline can propagate telemetry.
type ThumbnailError struct {
	SourceKey    string `json:"sourceKey"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (h *Handler) HandleRequest(ctx context.Context, event mediaprocessingmodel.StepFunctionPayload) (mediaprocessingmodel.LambdaResponse, error) {
	h.logger.Info("Thumbnails Lambda started", "job_id", event.JobID, "listing_identity_id", event.ListingIdentityID, "assets_to_process", len(event.Assets))

	allGeneratedAssets := make([]mediaprocessingmodel.JobAsset, 0)
	collectedErrors := make([]ThumbnailError, 0)
	bucket := os.Getenv("MEDIA_BUCKET")
	if bucket == "" {
		bucket = "toq-listing-medias"
	}

	for i, asset := range event.Assets {
		h.logger.Info("Inspecting asset", "index", i, "key", asset.Key, "type", asset.Type)

		if asset.Error != "" {
			h.logger.Warn("Skipping asset with previous error", "key", asset.Key, "error", asset.Error)
			continue
		}

		generatedKeys, err := h.service.ProcessImage(ctx, bucket, asset.Key)
		if err != nil {
			h.logger.Error("Failed to process asset", "key", asset.Key, "error", err)
			collectedErrors = append(collectedErrors, ThumbnailError{
				SourceKey:    asset.Key,
				ErrorCode:    "THUMBNAIL_PROCESSING_FAILED",
				ErrorMessage: err.Error(),
			})
			continue
		}

		for _, key := range generatedKeys {
			allGeneratedAssets = append(allGeneratedAssets, mediaprocessingmodel.JobAsset{
				Key:       key,
				Type:      asset.Type, // Or derived type?
				SourceKey: asset.Key,
			})
		}
	}

	h.logger.Info("Thumbnails Lambda finished", "generated_count", len(allGeneratedAssets))

	return mediaprocessingmodel.LambdaResponse{
		Body: ThumbnailOutput{
			Status:          "SUCCESS",
			GeneratedAssets: allGeneratedAssets,
			Errors:          collectedErrors,
		},
	}, nil
}
