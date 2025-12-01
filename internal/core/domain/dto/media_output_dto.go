package dto

// RequestUploadURLsOutput returns signed URLs ready to be used by the uploader.
type RequestUploadURLsOutput struct {
	ListingIdentityID   int64               `json:"listingIdentityId"`
	UploadURLTTLSeconds int                 `json:"uploadUrlTtlSeconds"`
	Files               []UploadInstruction `json:"files"`
}

// UploadInstruction carries the information required to perform a PUT upload to S3.
type UploadInstruction struct {
	AssetType string            `json:"assetType"`
	Sequence  uint8             `json:"sequence"`
	UploadURL string            `json:"uploadUrl"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	ObjectKey string            `json:"objectKey"`
	Title     string            `json:"title"`
}
