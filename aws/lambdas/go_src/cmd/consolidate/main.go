package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	consolidateservice "github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/service/consolidate"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// ConsolidateInput represents the combined input from Step Function
type ConsolidateInput struct {
	JobID             uint64                          `json:"jobId"`
	ListingIdentityID uint64                          `json:"listingIdentityId"`
	ExecutionArn      string                          `json:"executionArn"`
	StartedAt         string                          `json:"startedAt"`
	Assets            []mediaprocessingmodel.JobAsset `json:"assets"`
	ParallelResults   []ParallelResult                `json:"parallelResults"`
	Traceparent       string                          `json:"traceparent"`
}

// ParallelResult captures the generic output of parallel branches
type ParallelResult struct {
	Body struct {
		Thumbnails []mediaprocessingmodel.JobAsset `json:"generatedAssets"`
		Errors     []ThumbnailError                `json:"errors"`
		// Future: Videos []...
	} `json:"body"`
}

// ThumbnailError resumes failures reported by the thumbnail branch.
type ThumbnailError struct {
	SourceKey    string `json:"sourceKey"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func HandleRequest(ctx context.Context, event ConsolidateInput) (mediaprocessingmodel.LambdaResponse, error) {
	if event.JobID == 0 {
		return mediaprocessingmodel.LambdaResponse{}, fmt.Errorf("jobId is required")
	}
	if event.ListingIdentityID == 0 {
		return mediaprocessingmodel.LambdaResponse{}, fmt.Errorf("listingIdentityId is required")
	}

	// LOG: Full input (careful with size in prod, ok for debug)
	inputJSON, _ := json.Marshal(event)
	logger.Info("Consolidate Lambda started", "job_id", event.JobID, "listing_identity_id", event.ListingIdentityID, "raw_input_size", len(inputJSON))

	accumulators := consolidateservice.InitializePayloads(event.Assets)

	// 2. Process Parallel Results (Thumbnails)
	for idx, branch := range event.ParallelResults {
		thumbs := branch.Body.Thumbnails
		if len(thumbs) > 0 {
			logger.Info("Processing thumbnails results", "branch_index", idx, "count", len(thumbs))
		}

		for _, thumb := range thumbs {
			accumulator, exists := accumulators[thumb.SourceKey]
			if !exists {
				logger.Warn("Thumbnail without matching asset", "source_key", thumb.SourceKey, "thumb_key", thumb.Key)
				continue
			}

			consolidateservice.MapGeneratedAsset(accumulator, thumb)
		}

		if len(branch.Body.Errors) > 0 {
			logger.Warn("Thumbnail branch reported errors", "branch_index", idx, "count", len(branch.Body.Errors))
			consolidateservice.ApplyBranchErrors(accumulators, toBranchErrors(branch.Body.Errors))
		}
	}

	finalOutputs := consolidateservice.FlattenPayloads(accumulators)

	// LOG: Final output to be sent to backend
	outputJSON, _ := json.Marshal(finalOutputs)
	logger.Info("Consolidation complete",
		"job_id", event.JobID,
		"output_items", len(finalOutputs),
		"payload_preview", string(outputJSON),
	)

	responseBody := map[string]any{
		"jobId":             event.JobID,
		"listingIdentityId": event.ListingIdentityID,
		"executionArn":      event.ExecutionArn,
		"startedAt":         event.StartedAt,
		"provider":          string(mediaprocessingmodel.MediaProcessingProviderStepFunctions),
		"status":            string(mediaprocessingmodel.MediaProcessingJobStatusSucceeded),
		"failureReason":     "",
		"error":             nil,
		"outputs":           finalOutputs,
	}

	if event.Traceparent != "" {
		responseBody["traceparent"] = event.Traceparent
	}

	return mediaprocessingmodel.LambdaResponse{Body: responseBody}, nil
}

func toBranchErrors(errors []ThumbnailError) []consolidateservice.BranchError {
	converted := make([]consolidateservice.BranchError, 0, len(errors))
	for _, err := range errors {
		converted = append(converted, consolidateservice.BranchError{
			SourceKey:    err.SourceKey,
			ErrorCode:    err.ErrorCode,
			ErrorMessage: err.ErrorMessage,
		})
	}
	return converted
}

func main() {
	lambda.Start(HandleRequest)
}
