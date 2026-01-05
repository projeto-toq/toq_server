package videothumbnails

import (
	"context"
	"log/slog"
	"os"
	"strings"

	videoprocessing "github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/service/video_processing"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

// Handler orchestrates the video thumbnail Lambda execution flow.
type Handler struct {
	service *videoprocessing.VideoThumbnailService
	logger  *slog.Logger
}

// NewHandler builds the handler with its dependencies.
func NewHandler(service *videoprocessing.VideoThumbnailService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type Output struct {
	Status          string                          `json:"status"`
	GeneratedAssets []mediaprocessingmodel.JobAsset `json:"generatedAssets"`
	Errors          []ThumbnailError                `json:"errors"`
}

// ThumbnailError surfaces failures to the orchestrator for telemetry.
type ThumbnailError struct {
	SourceKey    string `json:"sourceKey"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (h *Handler) HandleRequest(ctx context.Context, event mediaprocessingmodel.StepFunctionPayload) (mediaprocessingmodel.LambdaResponse, error) {
	assets := event.VideoAssets
	if len(assets) == 0 {
		for _, asset := range event.Assets {
			if strings.Contains(strings.ToUpper(asset.Type), "VIDEO") {
				assets = append(assets, asset)
			}
		}
	}

	h.logger.Info("Video Thumbnails Lambda started",
		"job_id", event.JobID,
		"listing_identity_id", event.ListingIdentityID,
		"assets_to_process", len(assets),
	)

	bucket := resolveBucket()
	generated := make([]mediaprocessingmodel.JobAsset, 0, len(assets))
	errs := make([]ThumbnailError, 0)

	for i, asset := range assets {
		h.logger.Info("Processing video asset", "index", i, "key", asset.Key, "type", asset.Type)

		if asset.Error != "" {
			h.logger.Warn("Skipping asset with previous error", "key", asset.Key, "error", asset.Error)
			continue
		}

		outputKey, err := h.service.GenerateThumbnail(ctx, bucket, asset.Key)
		if err != nil {
			h.logger.Error("Video thumbnail generation failed", "key", asset.Key, "error", err)
			errs = append(errs, ThumbnailError{
				SourceKey:    asset.Key,
				ErrorCode:    "VIDEO_THUMBNAIL_FAILED",
				ErrorMessage: err.Error(),
			})
			continue
		}

		sourceKey := asset.SourceKey
		if sourceKey == "" {
			sourceKey = asset.Key
		}

		generated = append(generated, mediaprocessingmodel.JobAsset{
			Key:       outputKey,
			Type:      "VIDEO_THUMBNAIL",
			SourceKey: sourceKey,
		})
	}

	h.logger.Info("Video thumbnails completed", "generated_count", len(generated), "error_count", len(errs))

	return mediaprocessingmodel.LambdaResponse{Body: Output{
		Status:          "SUCCESS",
		GeneratedAssets: generated,
		Errors:          errs,
	}}, nil
}

func resolveBucket() string {
	bucket := os.Getenv("MEDIA_BUCKET")
	if bucket == "" {
		bucket = "toq-listing-medias"
	}
	return bucket
}
