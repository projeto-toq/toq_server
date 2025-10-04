package cnpjadapter

import (
	"fmt"
	"net/http"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

type CNPJAdapter struct {
	Client  *http.Client
	Token   string
	URLBase string
}

func NewCNPJAdapter(env *globalmodel.Environment) (*CNPJAdapter, error) {
	if env.CNPJ.Token == "" {
		return nil, fmt.Errorf("CNPJ token is required")
	}
	if env.CNPJ.URLBase == "" {
		return nil, fmt.Errorf("CNPJ URL base is required")
	}

	client := &http.Client{
		Timeout: 40 * time.Second,
	}
	return &CNPJAdapter{
		Client:  client,
		Token:   env.CNPJ.Token,
		URLBase: env.CNPJ.URLBase,
	}, nil
}
