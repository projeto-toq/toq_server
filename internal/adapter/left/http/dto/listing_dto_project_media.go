package dto

// RequestProjectUploadURLsRequest represents the payload to request signed URLs for project media uploads (plans/renders).
//
// This DTO is exclusive for OffPlanHouse listings while in StatusPendingPlanLoading. Only project asset types are allowed.
type RequestProjectUploadURLsRequest struct {
	// ListingIdentityID identifies the listing receiving the project media
	// Must correspond to a listing in StatusPendingPlanLoading
	// Example: 1024
	ListingIdentityID uint64 `json:"listingIdentityId" binding:"required,min=1" example:"1024"`

	// Files enumerates every project asset to upload
	// Must contain at least one file and only project asset types
	Files []ProjectUploadFileRequest `json:"files" binding:"required,min=1,dive"`
}

// ProjectUploadFileRequest describes a single project asset (plan PDF or render image) to be uploaded.
type ProjectUploadFileRequest struct {
	// AssetType must be PROJECT_DOC (PDF) or PROJECT_RENDER (image)
	// Example: "PROJECT_DOC"
	AssetType string `json:"assetType" binding:"required,oneof=PROJECT_DOC PROJECT_RENDER" enums:"PROJECT_DOC,PROJECT_RENDER" example:"PROJECT_DOC"`

	// Sequence determines the order of the asset within its type
	// Must be positive and unique per asset type
	// Example: 1
	Sequence uint8 `json:"sequence" binding:"required,min=1" example:"1"`

	// Filename is the original client filename
	// Example: "plantas_andar_1.pdf"
	Filename string `json:"filename" binding:"required" example:"plantas_andar_1.pdf"`

	// ContentType is the MIME type of the file
	// Allowed examples: application/pdf, image/jpeg, image/png
	ContentType string `json:"contentType" binding:"required" example:"application/pdf"`

	// Bytes represents the file size
	// Must be greater than zero and respect service limits
	// Example: 1048576
	Bytes int64 `json:"bytes" binding:"required,min=1" example:"1048576"`

	// Checksum is the SHA-256 hash of the file content (hex-encoded)
	// Example: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	Checksum string `json:"checksum" binding:"required" example:"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"`

	// Title is an optional caption for the asset
	// Example: "Planta pavimento térreo"
	Title string `json:"title,omitempty" example:"Planta pavimento térreo"`

	// Metadata carries optional key-value pairs for future processing hints
	Metadata map[string]string `json:"metadata,omitempty"`
}

// RequestProjectUploadURLsResponse returns the signed upload instructions for project media.
type RequestProjectUploadURLsResponse struct {
	// ListingIdentityID confirms the listing receiving the uploads
	ListingIdentityID uint64 `json:"listingIdentityId" example:"1024"`

	// UploadURLTTLSeconds indicates how long the signed URLs remain valid
	// Example: 900 (15 minutes)
	UploadURLTTLSeconds int `json:"uploadUrlTtlSeconds" example:"900"`

	// Files contains upload instructions for each project asset
	Files []RequestUploadInstructionResponse `json:"files"`
}

// CompleteProjectMediaRequest confirms all project media uploads and triggers finalization.
type CompleteProjectMediaRequest struct {
	// ListingIdentityID identifies the listing to finalize project media
	// Must be greater than zero
	// Example: 1024
	ListingIdentityID uint64 `json:"listingIdentityId" binding:"required,min=1" example:"1024"`
}
