package s3adapter

const (
	// Photo types - relative paths within user folder
	// These paths are used to construct full S3 object keys: /{userID}/{PhotoType}
	PhotoTypeOriginal = "photo.jpg"
	PhotoTypeSmall    = "thumbnails/small.jpg"
	PhotoTypeMedium   = "thumbnails/medium.jpg"
	PhotoTypeLarge    = "thumbnails/large.jpg"

	// CRECI Document types - relative paths within user folder
	// These paths are used to construct full S3 object keys: /{userID}/{DocumentType}
	CreciDocumentSelfie = "selfie.jpg"
	CreciDocumentFront  = "front.jpg"
	CreciDocumentBack   = "back.jpg"
)
