package mediaprocessingmodel

// MediaFinalizationInput defines the payload for the finalization pipeline.
type MediaFinalizationInput struct {
	JobID     uint64     `json:"jobId"`
	ListingID uint64     `json:"listingId"`
	Assets    []JobAsset `json:"assets"`
}
