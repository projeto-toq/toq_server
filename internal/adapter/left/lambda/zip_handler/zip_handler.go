package ziphandler

import (
	"context"

	"log/slog"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	lambdaport "github.com/projeto-toq/toq_server/internal/core/port/left/lambda_port"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

type ZipHandler struct {
	service lambdaport.ZipProcessingServiceInterface
	logger  *slog.Logger
}

func NewZipHandler(service lambdaport.ZipProcessingServiceInterface, logger *slog.Logger) *ZipHandler {
	return &ZipHandler{service: service, logger: logger}
}

// HandleRequest processes the Step Function event
// @Summary Generates ZIP bundles from processed assets
func (h *ZipHandler) HandleRequest(ctx context.Context, event mediaprocessingmodel.StepFunctionPayload) (mediaprocessingmodel.ZipOutput, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return mediaprocessingmodel.ZipOutput{}, err
	}
	defer spanEnd()

	h.logger.Info("Starting ZIP generation", "batch_id", event.BatchID, "listing_id", event.ListingID, "valid_assets_count", len(event.ValidAssets), "parallel_results_count", len(event.ParallelResults))

	// Extract thumbnails from ParallelResults (Fixing the root cause)
	thumbnails := h.extractThumbnails(event.ParallelResults)

	input := mediaprocessingmodel.GenerateZipInput{
		BatchID:     event.BatchID,
		ListingID:   event.ListingID,
		ValidAssets: event.ValidAssets,
		Thumbnails:  thumbnails,
	}

	output, err := h.service.GenerateZipBundle(ctx, input)
	if err != nil {
		h.logger.Error("Failed to generate ZIP", "err", err)
		return mediaprocessingmodel.ZipOutput{}, err
	}

	return output, nil
}

func (h *ZipHandler) extractThumbnails(results []mediaprocessingmodel.ParallelResult) []mediaprocessingmodel.MediaAssetDTO {
	var thumbnails []mediaprocessingmodel.MediaAssetDTO
	for _, result := range results {
		if len(result.Body.Thumbnails) > 0 {
			thumbnails = append(thumbnails, result.Body.Thumbnails...)
		}
	}
	return thumbnails
}
