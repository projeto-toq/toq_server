package cnpjadapter

import (
	"fmt"
	"net/http"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

type CNPJAdapter struct {
	Client   *http.Client
	Token    string
	URLBase  string
	Status   bool   `json:"status"`
	Return   string `json:"return"`
	Consumed int    `json:"consumed"`
	Result   struct {
		NumeroDeCNPJ   string `json:"numero_de_inscricao"`
		NomeDaPJ       string `json:"nome"`
		Fantasia       string `json:"fantasia"`
		DataNascimento string `json:"abertura"`
	} `json:"result"`
}

func NewCNPJAdapter(env *globalmodel.Environment) (*CNPJAdapter, error) {
	if env.CNPJ.Token == "" {
		return nil, fmt.Errorf("CNPJ token is required")
	}
	if env.CNPJ.URLBase == "" {
		return nil, fmt.Errorf("CNPJ URL base is required")
	}

	client := &http.Client{
		Timeout: 30 * time.Second, // Reduced from 600s for better UX
	}
	return &CNPJAdapter{
		Client:  client,
		Token:   env.CNPJ.Token,
		URLBase: env.CNPJ.URLBase,
	}, nil
}
