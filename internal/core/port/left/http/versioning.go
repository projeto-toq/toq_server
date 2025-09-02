package httpport

// APIVersionProvider expõe informações de versionamento/base path para a camada HTTP.
// Mantém a borda HTTP desacoplada de strings fixas e facilita futuras migrações.
type APIVersionProvider interface {
	// BasePath retorna o prefixo base da API, ex.: "/api/v2"
	BasePath() string
	// Version retorna a versão lógica, ex.: "v2"
	Version() string
}
