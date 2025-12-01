package zip

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/service/zip"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

type Handler struct {
	service *zip.ZipService
	logger  *slog.Logger
}

func NewHandler(service *zip.ZipService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type ZipOutput struct {
	ZipKey       string   `json:"zipKey"`
	AssetsZipped int      `json:"assetsZipped"`
	ZipBundles   []string `json:"zipBundles"`
}

func (h *Handler) HandleRequest(ctx context.Context, event mediaprocessingmodel.StepFunctionPayload) (mediaprocessingmodel.LambdaResponse, error) {
	h.logger.Info("Zip Lambda started", "job_id", event.JobID, "listing_id", event.ListingID, "assets_count", len(event.Assets))

	if len(event.Assets) == 0 {
		return mediaprocessingmodel.LambdaResponse{
			Body: ZipOutput{
				ZipKey:       "",
				AssetsZipped: 0,
				ZipBundles:   []string{},
			},
		}, nil
	}

	// Determine source keys
	var sourceKeys []string
	for _, asset := range event.Assets {
		sourceKeys = append(sourceKeys, asset.Key)
	}

	// Determine destination key
	bucket := os.Getenv("MEDIA_BUCKET")
	if bucket == "" {
		bucket = "toq-listing-medias"
	}

	zipKey := fmt.Sprintf("processed/zip/%d.zip", event.ListingID)

	err := h.service.CreateZip(ctx, bucket, sourceKeys, zipKey)
	if err != nil {
		h.logger.Error("Failed to process zip", "error", err)
		return mediaprocessingmodel.LambdaResponse{}, err
	}

	return mediaprocessingmodel.LambdaResponse{
		Body: ZipOutput{
			ZipKey:       zipKey,
			AssetsZipped: len(event.Assets),
			ZipBundles:   []string{zipKey},
		},
	}, nil
}
