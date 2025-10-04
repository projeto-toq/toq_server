package cpfadapter

import (
	"fmt"
	"net/http"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

type CPFAdapter struct {
	Client  *http.Client
	Token   string
	URLBase string
}

func NewCPFAdapter(env *globalmodel.Environment) (*CPFAdapter, error) {
	if env.CPF.Token == "" {
		return nil, fmt.Errorf("CPF token is required")
	}
	if env.CPF.URLBase == "" {
		return nil, fmt.Errorf("CPF URL base is required")
	}

	client := &http.Client{
		Timeout: 40 * time.Second,
	}
	return &CPFAdapter{
		Client:  client,
		Token:   env.CPF.Token,
		URLBase: env.CPF.URLBase,
	}, nil
}
