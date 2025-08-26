package s3adapter

const (
	// UsersBucketName é o nome do bucket principal para todos os tipos de media (users, listings, temp)
	UsersBucketName = "toq-app-media"

	// Photo types - caminhos relativos dentro da pasta do usuário
	PhotoTypeOriginal = "photo.jpg"
	PhotoTypeSmall    = "thumbnails/small.jpg"
	PhotoTypeMedium   = "thumbnails/medium.jpg"
	PhotoTypeLarge    = "thumbnails/large.jpg"

	// CRECI Document types - caminhos relativos dentro da pasta do usuário
	CreciDocumentSelfie = "selfie.jpg"
	CreciDocumentFront  = "front.jpg"
	CreciDocumentBack   = "back.jpg"
)
