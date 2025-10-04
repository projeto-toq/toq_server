package cepadapter

import (
	"fmt"
	"net/http"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

type CEPAdapter struct {
	Client  *http.Client
	Token   string
	URLBase string
}

func NewCEPAdapter(env *globalmodel.Environment) (*CEPAdapter, error) {
	if env.CEP.Token == "" {
		return nil, fmt.Errorf("CEP token is required")
	}
	if env.CEP.URLBase == "" {
		return nil, fmt.Errorf("CEP URL base is required")
	}

	client := &http.Client{
		Timeout: 30 * time.Second, // Reduced from 600s for better UX
	}
	return &CEPAdapter{
		Client:  client,
		Token:   env.CEP.Token,
		URLBase: env.CEP.URLBase,
	}, nil
}
