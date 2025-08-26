package storagemodel

// BucketConfig representa a configuração do bucket de armazenamento
type BucketConfig struct {
	Name   string
	Region string
}

// PhotoType representa os tipos de foto disponíveis
type PhotoType string

const (
	PhotoOriginal PhotoType = PhotoTypeOriginal
	PhotoSmall    PhotoType = PhotoTypeSmall
	PhotoMedium   PhotoType = PhotoTypeMedium
	PhotoLarge    PhotoType = PhotoTypeLarge
)

// DocumentType representa os tipos de documento disponíveis
type DocumentType string

const (
	DocSelfie DocumentType = CreciDocumentSelfie
	DocFront  DocumentType = CreciDocumentFront
	DocBack   DocumentType = CreciDocumentBack
)

// ValidPhotoTypes retorna um map com os tipos de foto válidos
func ValidPhotoTypes() map[string]bool {
	return map[string]bool{
		string(PhotoOriginal): true,
		string(PhotoSmall):    true,
		string(PhotoMedium):   true,
		string(PhotoLarge):    true,
	}
}

// ValidDocumentTypes retorna um map com os tipos de documento válidos
func ValidDocumentTypes() map[string]bool {
	return map[string]bool{
		string(DocSelfie): true,
		string(DocFront):  true,
		string(DocBack):   true,
	}
}
