package zip

import (
	"context"
	"errors"
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
	h.logger.Info("lambda.zip.start", "job_id", event.JobID, "listing_identity_id", event.ListingIdentityID, "assets_count", len(event.Assets))

	if err := h.validatePayload(event); err != nil {
		h.logger.Error("lambda.zip.invalid_payload", "error", err, "job_id", event.JobID)
		return mediaprocessingmodel.LambdaResponse{}, err
	}

	sourceKeys := extractAssetKeys(event.Assets)
	if len(sourceKeys) == 0 {
		h.logger.Warn("lambda.zip.no_assets", "job_id", event.JobID, "listing_identity_id", event.ListingIdentityID)
		return mediaprocessingmodel.LambdaResponse{
			Body: ZipOutput{ZipKey: "", AssetsZipped: 0, ZipBundles: []string{}},
		}, nil
	}

	bucket := h.resolveBucket()
	destinationKey := buildZipKey(event.ListingIdentityID, event.JobID)

	if err := h.service.CreateZip(ctx, bucket, sourceKeys, destinationKey); err != nil {
		h.logger.Error("lambda.zip.create_zip_error", "error", err, "bucket", bucket, "destination", destinationKey)
		return mediaprocessingmodel.LambdaResponse{}, err
	}

	h.logger.Info("lambda.zip.completed", "job_id", event.JobID, "listing_identity_id", event.ListingIdentityID, "zip_key", destinationKey)

	return mediaprocessingmodel.LambdaResponse{
		Body: ZipOutput{
			ZipKey:       destinationKey,
			AssetsZipped: len(sourceKeys),
			ZipBundles:   []string{destinationKey},
		},
	}, nil
}

func (h *Handler) validatePayload(event mediaprocessingmodel.StepFunctionPayload) error {
	if event.JobID == 0 {
		return errors.New("jobId is required")
	}
	if event.ListingIdentityID == 0 {
		return errors.New("listingIdentityId is required")
	}
	return nil
}

func (h *Handler) resolveBucket() string {
	bucket := os.Getenv("MEDIA_BUCKET")
	if bucket == "" {
		return "toq-listing-medias"
	}
	return bucket
}

func buildZipKey(listingID, jobID uint64) string {
	return fmt.Sprintf("%d/processed/zip/%d.zip", listingID, jobID)
}

func extractAssetKeys(assets []mediaprocessingmodel.JobAsset) []string {
	keys := make([]string, 0, len(assets))
	for _, asset := range assets {
		if asset.Key == "" {
			continue
		}
		keys = append(keys, asset.Key)
	}
	return keys
}
