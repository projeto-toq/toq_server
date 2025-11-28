package mediaprocessingmodel

// StepFunctionPayload represents the input from AWS Step Functions
type StepFunctionPayload struct {
	BatchID         string           `json:"batchId"`
	ListingID       uint64           `json:"listingId"`
	ValidAssets     []MediaAssetDTO  `json:"validAssets"`
	ParallelResults []ParallelResult `json:"parallelResults"` // Raw output from parallel state
}

type ParallelResult struct {
	Body struct {
		Thumbnails []MediaAssetDTO `json:"thumbnails"`
		Status     string          `json:"status"`
	} `json:"body"`
}

type MediaAssetDTO struct {
	SourceKey    string `json:"sourceKey"`
	ThumbnailKey string `json:"thumbnailKey"`
	AssetType    string `json:"assetType"`
}

type GenerateZipInput struct {
	BatchID     string
	ListingID   uint64
	ValidAssets []MediaAssetDTO
	Thumbnails  []MediaAssetDTO
}

type ZipOutput struct {
	ZipKey       string `json:"zipKey"`
	AssetsZipped int    `json:"assetsZipped"`
}
