package mediaprocessingservice

import mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"

// buildFinalizationInput translates processed assets to the payload expected by Step Functions.
func buildFinalizationInput(jobID uint64, listingID uint64, assets []mediaprocessingmodel.MediaAsset) mediaprocessingmodel.MediaFinalizationInput {
	jobAssets := make([]mediaprocessingmodel.JobAsset, 0, len(assets))
	for _, asset := range assets {
		jobAssets = append(jobAssets, mediaprocessingmodel.JobAsset{
			Key:  asset.S3KeyProcessed(),
			Type: string(asset.AssetType()),
		})
	}

	return mediaprocessingmodel.MediaFinalizationInput{
		JobID:     jobID,
		ListingID: listingID,
		Assets:    jobAssets,
	}
}
