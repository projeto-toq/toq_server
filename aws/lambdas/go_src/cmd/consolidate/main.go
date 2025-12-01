package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// ConsolidateInput represents the combined input from Step Function
type ConsolidateInput struct {
	JobID             string                          `json:"jobId"`
	ListingIdentityID uint64                          `json:"listingIdentityId"`
	Assets            []mediaprocessingmodel.JobAsset `json:"assets"`
	ParallelResults   []ParallelResult                `json:"parallelResults"`
}

// ParallelResult captures the generic output of parallel branches
type ParallelResult struct {
	Body struct {
		Thumbnails []mediaprocessingmodel.JobAsset `json:"generatedAssets"`
		// Future: Videos []...
	} `json:"body"`
}

func HandleRequest(ctx context.Context, event ConsolidateInput) (mediaprocessingmodel.LambdaResponse, error) {
	// LOG: Full input (careful with size in prod, ok for debug)
	inputJSON, _ := json.Marshal(event)
	logger.Info("Consolidate Lambda started", "job_id", event.JobID, "listing_identity_id", event.ListingIdentityID, "raw_input_size", len(inputJSON))

	// Map: SourceKey -> Output Payload
	resultsMap := make(map[string]*mediaprocessingmodel.MediaProcessingJobPayload)

	// 1. Initialize with original assets
	for _, asset := range event.Assets {
		resultsMap[asset.Key] = &mediaprocessingmodel.MediaProcessingJobPayload{
			RawKey:  asset.Key, // The Backend looks for THIS
			Outputs: make(map[string]string),
		}
	}

	// 2. Process Parallel Results (Thumbnails)
	if len(event.ParallelResults) > 0 {
		thumbs := event.ParallelResults[0].Body.Thumbnails
		logger.Info("Processing thumbnails results", "count", len(thumbs))

		for _, thumb := range thumbs {
			if payload, exists := resultsMap[thumb.SourceKey]; exists {
				// Map thumbnail
				payload.Outputs["thumbnail_"+thumb.Type] = thumb.Key

				// Define main thumbnail (e.g., MEDIUM)
				if thumb.Type == "THUMBNAIL_MEDIUM" {
					payload.ThumbnailKey = thumb.Key
				}

				logger.Debug("Mapped thumbnail", "source", thumb.SourceKey, "thumb", thumb.Key)
			} else {
				logger.Warn("Orphaned thumbnail found", "source_key", thumb.SourceKey, "thumb_key", thumb.Key)
			}
		}
	}

	// 3. Convert to final list
	finalOutputs := make([]mediaprocessingmodel.MediaProcessingJobPayload, 0, len(resultsMap))
	for _, payload := range resultsMap {
		finalOutputs = append(finalOutputs, *payload)
	}

	// LOG: Final output to be sent to backend
	outputJSON, _ := json.Marshal(finalOutputs)
	logger.Info("Consolidation complete",
		"job_id", event.JobID,
		"output_items", len(finalOutputs),
		"payload_preview", string(outputJSON),
	)

	return mediaprocessingmodel.LambdaResponse{
		Body: map[string]any{
			"jobId":             event.JobID,
			"listingIdentityId": event.ListingIdentityID,
			"outputs":           finalOutputs,
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
