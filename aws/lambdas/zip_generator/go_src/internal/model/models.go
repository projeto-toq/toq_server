package model

// StepFunctionPayload represents the input from AWS Step Functions
type StepFunctionPayload struct {
	BatchID         uint64           `json:"batchId"`
	ListingID       uint64           `json:"listingId"`
	ValidAssets     []MediaAssetDTO  `json:"validAssets"`
	ParallelResults []ParallelResult `json:"parallelResults"`
}

type ParallelResult struct {
	Body struct {
		Status     string          `json:"status"`
		Thumbnails []MediaAssetDTO `json:"thumbnails"`
	} `json:"body"`
}

type MediaAssetDTO struct {
	SourceKey    string `json:"sourceKey"`
	ThumbnailKey string `json:"thumbnailKey"`
	AssetType    string `json:"assetType"`
}

type ZipOutput struct {
	ZipKey       string `json:"zipKey"`
	AssetsZipped int    `json:"assetsZipped"`
}
