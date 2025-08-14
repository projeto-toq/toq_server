package gcsadapter

const (
	// UsersBucketName é o nome do bucket único para todos os usuários
	UsersBucketName = "toq_server_users_media"

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
