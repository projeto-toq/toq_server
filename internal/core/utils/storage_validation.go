package utils

// IsAllowedImageContentType verifica se o content-type é permitido para upload de imagens.
// Somente image/jpeg e image/png são aceitos no momento.
func IsAllowedImageContentType(ct string) bool {
	switch ct {
	case "image/jpeg", "image/png":
		return true
	default:
		return false
	}
}
