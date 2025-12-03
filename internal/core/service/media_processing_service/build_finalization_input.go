package mediaprocessingservice

import (
	"context"
	"fmt"

	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"go.opentelemetry.io/otel/trace"
)

// buildFinalizationInput translates processed assets to the payload expected by Step Functions.
func buildFinalizationInput(ctx context.Context, jobID uint64, listingID uint64, assets []mediaprocessingmodel.MediaAsset) mediaprocessingmodel.MediaFinalizationInput {
	jobAssets := make([]mediaprocessingmodel.JobAsset, 0, len(assets))
	for _, asset := range assets {
		jobAssets = append(jobAssets, mediaprocessingmodel.JobAsset{
			Key:  asset.S3KeyProcessed(),
			Type: string(asset.AssetType()),
		})
	}

	return mediaprocessingmodel.MediaFinalizationInput{
		JobID:             jobID,
		ListingIdentityID: listingID,
		Assets:            jobAssets,
		Traceparent:       traceparentFromContext(ctx),
	}
}

func traceparentFromContext(ctx context.Context) string {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if !spanCtx.IsValid() {
		return ""
	}
	return fmt.Sprintf("00-%s-%s-%s", spanCtx.TraceID().String(), spanCtx.SpanID().String(), spanCtx.TraceFlags().String())
}
