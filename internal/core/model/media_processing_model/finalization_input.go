package mediaprocessingmodel

// MediaFinalizationInput defines the payload for the finalization pipeline.
type MediaFinalizationInput struct {
	JobID             uint64     `json:"jobId"`
	ListingIdentityID uint64     `json:"listingIdentityId"`
	Assets            []JobAsset `json:"assets"`
	Traceparent       string     `json:"traceparent,omitempty"`
}
