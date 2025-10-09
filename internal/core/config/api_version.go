package config

import httpport "github.com/projeto-toq/toq_server/internal/core/port/left/http"

// staticAPIVersionProvider é uma implementação simples que retorna a versão atual
// da API. Mantém a string centralizada e facilita futuras evoluções.
type staticAPIVersionProvider struct{}

func NewStaticAPIVersionProvider() httpport.APIVersionProvider {
	return &staticAPIVersionProvider{}
}

func (p *staticAPIVersionProvider) BasePath() string { return "/api/v2" }
func (p *staticAPIVersionProvider) Version() string  { return "v2" }
